package database

import (
	"database/sql"
	"github.com/mvl-at/qbs"
	"rest/context"
)

func log(err error) {

	if err != nil {
		context.ErrLog.Println(err)
	}
}

func TableCreate(a interface{}) {
	m, _ := qbs.GetMigration()

	err := m.CreateTableIfNotExists(a)
	log(err)
}

func Delete(a interface{}) {
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	_, err = db.Delete(a)
	if err != nil {
		log(err)
	}
}

func Save(a interface{}) {
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	_, err = db.Save(a)
	if err != nil {
		log(err)
	}
}

func FindAll(a interface{}) {
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	err = db.FindAll(a)
	if err != nil {
		log(err)
	}
}

func FindAllWhereEqual(a interface{}, field string, value interface{}) {
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	err = db.WhereEqual(field, value).FindAll(a)
	if err != nil {
		log(err)
	}
}

func Find(a interface{}) (exists bool) {
	exists = false
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	err = db.Find(a)

	if err != nil {

		if err != sql.ErrNoRows {
			log(err)
		}
	} else {
		exists = true
	}
	return
}
