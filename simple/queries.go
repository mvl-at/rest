package simple

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/rest/context"
)

var db *sql.DB

func OpenDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", context.Conf.SQLiteFile)
	if err != nil {
		return err
	}
	return err
}

func QueryData(query string, data *[]DBO, numFields int) error {
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for counter := 0; rows.Next(); counter++ {
		dbo, err := (*data)[0].Scan(rows.Scan, data)
		if err != nil {
			return err
		}
		if dbo != nil && counter > 0 {
			*data = append(*data, *dbo)
		}
		if counter == 0 {
			(*data)[0] = *dbo
		}
	}
	return nil
}
