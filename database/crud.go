package database

import (
	"github.com/mvl-at/qbs"
	"rest/context"
)

func log(err error) {

	if err != nil {
		context.ErrLog.Println(err)
	}
}

func GenericCreate(a interface{}) {
	m, _ := qbs.GetMigration()

	err := m.CreateTableIfNotExists(a)
	log(err)
}

func GenericDelete(a interface{}) {
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

func GenericSave(a interface{}) {
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

func GenericFetch(a interface{}) {
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

func GenericFetchWhereEqual(a interface{}, field string, value interface{}) {
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

func GenericSingleFetch(a interface{}) {
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	err = db.Find(a)
	if err != nil {
		log(err)
	}
}
