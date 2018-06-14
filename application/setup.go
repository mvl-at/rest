package application

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/qbs"
	"rest/context"
	"rest/database"
	"rest/http"
)

func Setup() {
	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	database.CheckRoot()
	http.Routes()
	http.Run()
}
