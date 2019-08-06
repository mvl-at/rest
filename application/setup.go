package application

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/qbs"
	"github.com/mvl-at/rest/context"
	"github.com/mvl-at/rest/database"
	"github.com/mvl-at/rest/http"
	"github.com/mvl-at/rest/simple"
	"math/rand"
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
	simple.OpenDatabase()
	simple.PersistenceRunner()
	http.Run()
}
