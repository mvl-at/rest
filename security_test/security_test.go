package security_test

import (
	_ "github.com/mattn/go-sqlite3"
	"rest/model"
	"testing"
)

func TestPaulInsertTuba(t *testing.T) {
	tuba := &model.Instrument{Name: "Tuba", NamePlural: "Tuben"}
	saveData(tuba, true, willi, t)
}
