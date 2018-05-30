package database

import "github.com/coocood/qbs"

func GenericCreate(a interface{}) {
	m, _ := qbs.GetMigration()
	m.CreateTableIfNotExists(a)
}

func GenericDelete(a interface{}) {
	db.Delete(a)
}

func GenericSave(a interface{}, autoKey bool) {
	db.Save(a)
}

func GenericFetch(a interface{}) {
	db.FindAll(a)
}
