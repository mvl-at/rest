package security_test

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/model"
	"github.com/mvl-at/rest/database"
	. "github.com/mvl-at/rest/mock"
	"testing"
	"time"
)

func TestPaulInsertTuba(t *testing.T) {
	tuba := &model.Instrument{Name: "Tuba", NamePlural: "Tuben"}
	saveData(tuba, true, Paulina, t)
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
	newJosef := &model.Member{Id: Josef.Id}
	database.Find(newJosef)
	deleteData(newJosef, true, Josef, t)
}

func TestKarlUpdateOwnPassword(t *testing.T) {
	credentials := &database.Credentials{MemberId: Karl.Id, Password: Karl.Password + "67", Username: Karl.Username + "df"}
	updateCredentials(credentials, false, true, Karl, t)
}

func TestKarlUpdateHelmutCredentials(t *testing.T) {
	credentials := &database.Credentials{MemberId: Helmut.Id, Password: Helmut.Password + "67", Username: Karl.Username + "df"}
	updateCredentials(credentials, false, false, Karl, t)
}

func TestJosefUpdateHelmutCredentials(t *testing.T) {
	credentials := &database.Credentials{MemberId: Helmut.Id, Password: Helmut.Password + "67", Username: Helmut.Username + "df"}
	updateCredentials(credentials, true, true, Josef, t)
}

func TestFranzUpdateKarlCredentials(t *testing.T) {
	credentials := &database.Credentials{MemberId: Karl.Id, Password: Karl.Password + "67", Username: Karl.Username + "df"}
	updateCredentials(credentials, false, false, Franz, t)
}

func TestFranzUpdatePaulCredentials(t *testing.T) {
	credentials := &database.Credentials{MemberId: Paulina.Id, Password: Paulina.Password + "67", Username: Paulina.Username}
	updateCredentials(credentials, false, true, Franz, t)
}

func TestPaulinaDeleteJosef(t *testing.T) {
	members := make([]*model.Member, 0)
	database.FindAll(&members)
	fmt.Println(members)
	deleteData(Josef, true, Paulina, t)
	members = make([]*model.Member, 0)
	database.FindAll(&members)
	for _, v := range members {
		fmt.Println(*v)
	}
}

func TestJosefDeletePaulina(t *testing.T) {
	deleteData(&model.Member{Id: Paulina.Id}, false, Josef, t)
	members := make([]*model.Member, 0)
	database.FindAll(&members)
	for _, v := range members {
		fmt.Println(*v)
	}
}

func TestJosefChangesNameButNotPassword(t *testing.T) {
	Josef.LastName = "Strauchi"
	saveData(Josef, true, Josef, t)
	nJosef := &model.Member{Id: Josef.Id}
	database.Find(nJosef)
	if nJosef.Password == "" {
		t.Errorf("password of %s is empty but should not", Josef.FirstName)
	}

}
