package database

import (
	"database/sql"
	"github.com/mvl-at/qbs"
	"reflect"
	"rest/context"
	"rest/model"
)

//Generic error for all database related actions.
type DBError struct {
	message string
}

func (d DBError) Error() string {
	return d.message
}

//Shortcut for error logging.
func log(err error) {

	if err != nil {
		context.ErrLog.Println(err)
	}
}

//Creates tables for all registered structs, if not exist.
func TableCreate(a interface{}) {
	m, _ := qbs.GetMigration()

	err := m.CreateTableIfNotExists(a)
	log(err)
}

//Deletes a single struct.
func Delete(a interface{}) {
	db, err := qbs.GetQbs()
	defer db.Close()
	db.Log = true

	if err != nil {
		log(err)
		return
	}

	if reflect.TypeOf(a) == reflect.TypeOf(&model.Member{}) {
		member := a.(*model.Member)
		roleMemberRoot := make([]*model.RoleMember, 0)
		db.WhereEqual("role_id", "root").FindAll(&roleMemberRoot)
		lastRoot := len(roleMemberRoot) <= 1

		if lastRoot && roleMemberRoot[0].MemberId == member.Id {
			log(DBError{"cannot delete last root!"})
			return
		}
	}

	_, err = db.Delete(a)
	if err != nil {
		log(err)
	}
}

//Persists a single struct.
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

//Finds a slice of structs.
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

//Finds a slice of structs with given condition.
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

//Finds a single struct with an Id.
//Returns false if no struct exists with the given Id.
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
