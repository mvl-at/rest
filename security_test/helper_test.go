package security_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mvl-at/model"
	"github.com/mvl-at/qbs"
	"github.com/mvl-at/rest/context"
	"github.com/mvl-at/rest/database"
	"github.com/mvl-at/rest/http"
	"github.com/mvl-at/rest/mock"
	"github.com/mvl-at/rest/security"
	"io/ioutil"
	"math/rand"
	vhttp "net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var willi = &model.Member{Username: "willi", Password: "123456"}

func init() {
	os.Remove(context.Conf.SQLiteFile)

	config := &context.Configuration{
		JwtExpiration: 100,
		JwtSecret:     string(rand.Int()),
		SQLiteFile:    "mvl.sqlite",
		Host:          "127.0.0.1",
		Port:          17534}

	context.Conf = config
	context.Log.SetFlags(0)
	context.Log.SetOutput(ioutil.Discard)
	context.ErrLog.SetFlags(0)
	context.ErrLog.SetOutput(ioutil.Discard)

	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	mock.MockData()
	http.Routes()
	go http.Run()
}

func pathName(a interface{}) (path string) {
	aType := reflect.TypeOf(a).Elem().Name()
	switch aType {
	case reflect.TypeOf(model.Instrument{}).Name():
		path = "instruments"
	case reflect.TypeOf(model.Member{}).Name():
		path = "members"
	case reflect.TypeOf(model.Role{}).Name():
		path = "roles"
	case reflect.TypeOf(model.Event{}).Name():
		path = "events"
	case reflect.TypeOf(model.LeaderRole{}).Name():
		path = "leaderRoles"
	case reflect.TypeOf(model.RoleMember{}).Name():
		path = "leaderRolesMembers"
	case reflect.TypeOf(model.LeaderRoleMember{}).Name():
		path = "leaderRolesMembers"
	}
	return
}

func saveData(a interface{}, shouldExistAfter bool, issuer *model.Member, t *testing.T) {
	request("/"+pathName(a), vhttp.MethodPost, a, issuer)
	if exists(a) != shouldExistAfter {

		if shouldExistAfter {
			t.Errorf("%s should exist but does not!", a)
		} else {
			t.Errorf("%s should not exist but does!", a)
		}
	}
}

func deleteData(a interface{}, shouldExistAfter bool, issuer *model.Member, t *testing.T) {
	request("/"+pathName(a), vhttp.MethodDelete, a, issuer)
	if exists(a) != shouldExistAfter {

		if shouldExistAfter {
			t.Errorf("%s should exist but does not!", a)
		} else {
			t.Errorf("%s should not exist but does!", a)
		}
	}
}

func exists(equality interface{}) bool {
	sli := reflect.New(reflect.SliceOf(reflect.TypeOf(equality))).Interface()
	database.FindAll(sli)

	sliceValue := reflect.ValueOf(sli).Elem()

	for i := 0; i < sliceValue.Len(); i++ {
		if equal(equality, sliceValue.Index(i).Interface()) {
			return true
		}
	}
	return false
}

func token(member *model.Member) string {
	data := &security.JWTData{Username: member.Username, Password: member.Password}
	jsonData, _ := json.Marshal(data)
	req, _ := vhttp.NewRequest(vhttp.MethodGet, fmt.Sprintf("http://%s:%d/login", context.Conf.Host, context.Conf.Port), bytes.NewBuffer(jsonData))
	c := &vhttp.Client{}
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func request(url string, method string, data interface{}, issuer *model.Member) (response string, status int) {
	var jsonData []byte
	if data != nil {
		jsonData, _ = json.Marshal(data)
	}
	req, _ := vhttp.NewRequest(method, fmt.Sprintf("http://%s:%d%s", context.Conf.Host, context.Conf.Port, url), bytes.NewBuffer(jsonData))

	if issuer != nil {
		req.Header.Set("token", token(issuer))
	}

	c := &vhttp.Client{}
	resp, _ := c.Do(req)

	if resp.Body != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		response = string(body)
	}
	status = resp.StatusCode
	return
}

func equal(a interface{}, b interface{}) bool {
	aType := reflect.TypeOf(a).Elem()
	if aType.Name() != reflect.TypeOf(b).Elem().Name() {
		return false
	}
	aValue := reflect.ValueOf(a).Elem()
	bValue := reflect.ValueOf(b).Elem()

	equal := true

	for i := 0; i < aValue.NumField() && equal; i++ {
		aField := aValue.Field(i)
		bField := bValue.Field(i)

		fieldName := aType.Field(i).Name

		if !strings.HasSuffix(fieldName, "Id") && reflect.TypeOf(aField.Interface()).Name() != reflect.TypeOf(time.Time{}).Name() && aField.Kind() != reflect.Ptr {
			equal = aField.Interface() == bField.Interface()
		} else {
			aTime, ok := aField.Interface().(time.Time)

			if ok {
				bTime, _ := bField.Interface().(time.Time)
				equal = timeEqual(aTime, bTime)
			}
		}

	}
	return equal
}

func timeEqual(t1 time.Time, t2 time.Time) bool {
	return t1.Minute() == t2.Minute() && t1.Hour() == t2.Hour() && t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year() && t1.Second() == t2.Second()
}

func updateCredentials(credentials *security.Credentials, username bool, password bool, issuer *model.Member, t *testing.T) {
	oldMember := &model.Member{Id: credentials.MemberId}
	database.Find(oldMember)
	_, status := request("/credentials", vhttp.MethodPost, credentials, issuer)
	if status == vhttp.StatusForbidden {
		if username && password {
			t.Errorf("%+v should be able to set credentials %+v but was not able to!", issuer, credentials)
		}
		return
	}
	newMember := &model.Member{Id: credentials.MemberId}
	database.Find(newMember)

	if username && oldMember.Username == newMember.Username {
		t.Errorf("%+v should be able to set username of %+v but was not able to!", issuer, credentials)
	}

	if password && oldMember.Password == newMember.Password {
		t.Errorf("%+v should be able to set password of %+v but was not able to!", issuer, credentials)
	}

	if !username && oldMember.Username != newMember.Username {
		t.Errorf("%+v should not be able to set username of %+v but was able to!", issuer, credentials)
	}

	if !password && oldMember.Password != newMember.Password {
		t.Errorf("%+v should not be able to set password of %+v but was able to!", issuer, credentials)
	}
}
