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

type instrument struct {
	model.Instrument
}

type event struct {
	model.Event
}

type member struct {
	model.Member
}

type role struct {
	model.Role
}

type leaderRole struct {
	model.LeaderRole
}

type roleMember struct {
	model.RoleMember
}

type leaderRoleMember struct {
	model.LeaderRoleMember
}

type testEquality interface {
	equal(other testEquality) bool
}

func (i instrument) equal(other testEquality) bool {
	otherIns, ok := other.(instrument)

	if !ok {
		return false
	}

	return otherIns.Name == i.Name && otherIns.NamePlural == i.NamePlural
}

func (e event) equal(other testEquality) bool {
	o, ok := other.(event)

	if !ok {
		return false
	}

	return e.Name == o.Name &&
		e.Internal == o.Internal &&
		e.Important == o.Important &&
		e.Note == o.Note &&
		e.Date.Unix() == o.Date.Unix() &&
		e.MusicianTime.Unix() == o.Date.Unix() &&
		e.Time.Unix() == o.Time.Unix() &&
		e.Place == o.Place &&
		e.MusicianPlace == o.MusicianPlace &&
		e.Uniform == o.Uniform
}

func (m member) equal(other testEquality) bool {
	o, ok := other.(member)

	if !ok {
		return false
	}

	return m.InstrumentId == o.InstrumentId &&
		m.LoginAllowed == o.LoginAllowed &&
		m.Username == o.Username &&
		m.Active == o.Active &&
		m.Deleted == o.Deleted &&
		m.Joined == o.Joined &&
		m.LastName == o.LastName &&
		m.FirstName == o.FirstName &&
		m.Picture == o.Picture
}

func (r role) equal(other testEquality) bool {
	o, ok := other.(role)

	if !ok {
		return false
	}

	return r.Name == o.Name &&
		r.NamePlural == o.NamePlural
}

func (l leaderRole) equal(other testEquality) bool {
	o, ok := other.(leaderRole)

	if !ok {
		return false
	}
	return l.NamePlural == o.NamePlural &&
		l.Name == o.Name
}

func (r roleMember) equal(other testEquality) bool {
	o, ok := other.(roleMember)

	if !ok {
		return false
	}

	return r.RoleId == o.RoleId &&
		r.MemberId == o.MemberId
}

func (l leaderRoleMember) equal(other testEquality) bool {
	o, ok := other.(leaderRoleMember)

	if !ok {
		return false
	}

	return l.LeaderRoleId == o.LeaderRoleId &&
		l.MemberId == o.MemberId &&
		l.Priority == o.Priority
}
