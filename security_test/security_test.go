package security_test

import (
	"bytes"
	"encoding/json"
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

}

func token(member *model.Member) string {
	data := &security.JWTData{Username: member.Username, Password: member.Password}
	jsonData, _ := json.Marshal(data)
	req, _ := vhttp.NewRequest(vhttp.MethodGet, "127.0.0.1:8080/login", bytes.NewBuffer(jsonData))
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
	req, _ := vhttp.NewRequest(method, "127.0.0.1:8080"+url, bytes.NewBuffer(jsonData))
	c := &vhttp.Client{}
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.StatusCode
}
