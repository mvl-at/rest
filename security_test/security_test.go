package security_test

import (
	_ "github.com/mattn/go-sqlite3"
	vhttp "net/http"
	"rest/database"
	"rest/model"
	"testing"
)

func TestInsert(t *testing.T) {
	tuba := &model.Instrument{Name: "Tuba", NamePlural: "Tuben"}
	request("/instruments", vhttp.MethodPost, tuba, willi)
	tuben := make([]*model.Instrument, 0)
	database.FindAll(&tuben)

	correct := false

	for _, v := range tuben {
		if v.NamePlural == tuba.NamePlural && v.Name == tuba.Name {
			correct = true
			break
		}
	}

	if !correct {
		t.Error("tuba was not inserted but should")
	}
}
