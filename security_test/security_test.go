package security_test

import (
	"github.com/mvl-at/qbs"
	vhttp "net/http"
	"os"
	"rest/context"
	"rest/database"
	"rest/http"
	"rest/mock"
	"testing"
)

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

	vhttp.hea
}
