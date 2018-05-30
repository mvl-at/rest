package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func Connect() {
	database, err := sql.Open("sqlite3", "mvl.sqlite")

	if err != nil {
		log.Fatal(err.Error())
	}
	db = database
}
