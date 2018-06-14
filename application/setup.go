package application

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/qbs"
	"math/rand"
	"rest/context"
	"rest/database"
	"rest/http"
	"time"
)

func Setup() {
	rand.Seed(time.Now().UnixNano())
	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	database.CheckRoot()
	http.Routes()
	http.Run()
}
