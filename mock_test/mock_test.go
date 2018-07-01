package mock_test

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/qbs"
	"github.com/mvl-at/rest/context"
	"github.com/mvl-at/rest/database"
	"github.com/mvl-at/rest/http"
	"github.com/mvl-at/rest/mock"
	"os"
	"testing"
)

func TestRunMock(t *testing.T) {
	os.Remove(context.Conf.SQLiteFile)
	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	mock.MockData()
	http.Routes()
	http.Run()
}
