package database

import (
	"github.com/coocood/qbs"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *qbs.Qbs

func Connect() {
	qbs.RegisterSqlite3("mvl.sqlite")
	database, err := qbs.GetQbs()

	if err != nil {
		log.Fatal(err.Error())
	}
	db = database
}
