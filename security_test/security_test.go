package security_test

import (
	_ "github.com/mattn/go-sqlite3"
	"rest/database"
	. "rest/mock"
	"rest/model"
	"testing"
	"time"
)

func TestPaulInsertTuba(t *testing.T) {
	tuba := &model.Instrument{Name: "Tuba", NamePlural: "Tuben"}
	saveData(tuba, true, Paul, t)
}

func TestJosefInsertTrombone(t *testing.T) {
	trombone := &model.Instrument{Name: "Posaune", NamePlural: "Posaunen"}
	saveData(trombone, true, Josef, t)
}

func TestHelmutInsertTrumpet(t *testing.T) {
	trumpet := &model.Instrument{Name: "Trompete", NamePlural: "Trompete"}
	saveData(trumpet, false, Helmut, t)
}

func TestJosefInsertKrampuskonzert(t *testing.T) {
	event := &model.Event{Name: "Krampuskonzert", Time: time.Date(1, 1, 1, 19, 0, 0, 0, time.Local),
		Uniform:       "Uniform ohne Hut",
		MusicianPlace: "Gasthaus List",
		Place:         "Gasthaus List",
		Date:          time.Date(2018, 12, 5, 0, 0, 0, 0, time.Local),
		MusicianTime:  time.Date(1, 1, 1, 18, 0, 0, 0, time.Local),
		Note:          "Anschl. Essen",
		Important:     true,
		Internal:      false}
	saveData(event, true, Josef, t)
}

func TestFranzInsertOktoberfest(t *testing.T) {
	event := &model.Event{Name: "Oktoberfest", Time: time.Date(1, 1, 1, 19, 0, 0, 0, time.Local),
		Uniform:       "Lederhose mit Polo",
		MusicianPlace: "Gutshof Prosoroff",
		Place:         "Gutshof Prosoroff",
		Date:          time.Date(2018, 9, 5, 0, 0, 0, 0, time.Local),
		MusicianTime:  time.Date(1, 1, 1, 18, 0, 0, 0, time.Local),
		Note:          "Anschl. Essen",
		Important:     true,
		Internal:      false}
	saveData(event, false, Franz, t)
}

func TestHelmutInsertWeihnachtsfeier(t *testing.T) {
	event := &model.Event{Name: "Weihnachtsfeier", Time: time.Date(1, 1, 1, 19, 0, 0, 0, time.Local),
		Uniform:       "Abendkleidung",
		MusicianPlace: "Musikheim",
		Date:          time.Date(2018, 12, 12, 0, 0, 0, 0, time.Local),
		MusicianTime:  time.Date(1, 1, 1, 18, 0, 0, 0, time.Local),
		Note:          "Mit Begleitung",
		Important:     false,
		Internal:      true}
	saveData(event, true, Helmut, t)
}

func TestKeepLastRoot(t *testing.T) {
	members := make([]*model.Member, 0)
	database.FindAll(&members)
	for _, member := range members {
		database.Delete(member)
	}
	members = make([]*model.Member, 0)
	database.FindAll(&members)
	rolesMembers := make([]*model.RoleMember, 0)
	database.FindAll(&rolesMembers)
	someOneWithRoot := false
	for _, member := range members {
		for _, role := range rolesMembers {
			if role.MemberId == member.Id && role.RoleId == "root" {
				someOneWithRoot = true
			}
		}
	}
	if !someOneWithRoot {
		t.Error("all members with root are deleted!")
	}
}
