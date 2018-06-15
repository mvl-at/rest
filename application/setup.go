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

//Does the setup for the application. That includes: initialize random, set logging tools,
//register and setup database, create root user when not existent, register http routes and run http
func Setup() {
	rand.Seed(time.Now().UnixNano())
	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	database.CheckRoot()
	http.Routes()
	http.Run()
}
