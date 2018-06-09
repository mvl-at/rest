package security_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mvl-at/qbs"
	"io/ioutil"
	"math/rand"
	vhttp "net/http"
	"os"
	"reflect"
	"rest/context"
	"rest/database"
	"rest/http"
	"rest/mock"
	"rest/model"
	"rest/security"
	"testing"
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

func saveData(a interface{}, issuer *model.Member, t *testing.T)   {}
func deleteData(a interface{}, issuer *model.Member, t *testing.T) {}

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
