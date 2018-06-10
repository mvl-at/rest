package security_test

import (
	_ "github.com/mattn/go-sqlite3"
	. "rest/mock"
	"rest/model"
	"testing"
)

func TestPaulInsertTuba(t *testing.T) {
	tuba := &model.Instrument{Name: "Tuba", NamePlural: "Tuben"}
	saveData(tuba, true, Paul, t)
}

func TestJosefInsertTrombone(t *testing.T) {
	trombone := &model.Instrument{Name: "Posaune", NamePlural: "Posaunen"}
	saveData(trombone, true, Josef, t)
}
