package security_test

import (
	"bytes"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/qbs"
	"io/ioutil"
	vhttp "net/http"
	"os"
	"rest/context"
	"rest/database"
	"rest/http"
	"rest/mock"
	"rest/model"
	"rest/security"
	"testing"
)

var willi = &model.Member{Username: "willi", Password: "123456"}

func setup() {
	os.Remove(context.Conf.SQLiteFile)
	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	mock.MockData()
	http.Routes()
	go http.Run()
}

func TestInsert(t *testing.T) {
	setup()
	tuba := &model.Instrument{Name: "Tuba", NamePlural: "Tuben"}
	request("/instruments", vhttp.MethodPost, tuba, willi)
	tuben := make([]*model.Instrument, 0)
	database.FindAll(&tuben)

	correct := false

	for _, v := range tuben {
		if v.NamePlural == tuba.NamePlural && v.Name == tuba.Name {
			correct = true
			break
		}
	}

	if !correct {
		t.Error("tuba was not inserted but should")
	}
}

func token(member *model.Member) string {
	data := &security.JWTData{Username: member.Username, Password: member.Password}
	jsonData, _ := json.Marshal(data)
	req, _ := vhttp.NewRequest(vhttp.MethodGet, "http://127.0.0.1:8080/login", bytes.NewBuffer(jsonData))
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
	req, _ := vhttp.NewRequest(method, "http://127.0.0.1:8080"+url, bytes.NewBuffer(jsonData))

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
