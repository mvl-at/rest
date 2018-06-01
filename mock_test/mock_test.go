package mock_test

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvl-at/qbs"
	"os"
	"rest/context"
	"rest/database"
	"rest/http"
	"rest/model"
	"testing"
	"time"
)

func TestRunMock(t *testing.T) {
	os.Remove(context.Conf.SQLiteFile)
	qbs.SetLogger(context.Log, context.ErrLog)
	qbs.RegisterSqlite3(context.Conf.SQLiteFile)
	database.Register()
	mockData()
	http.Routes()
	http.Run()
}

func mockData() {
	events := []*model.Event{
		{
			Name:          "Frühschoppen",
			Date:          time.Date(2018, 6, 17, 0, 0, 0, 0, time.Local),
			Uniform:       "MVL Polo und Lederhose",
			MusicianPlace: "Musikheim",
			Place:         "Musikheim",
			Time:          time.Date(1, 1, 1, 10, 0, 0, 0, time.Local),
			MusicianTime:  time.Date(0, 0, 0, 8, 45, 0, 0, time.Local),
			Note:          "Anschl. Wegräumen",
			Important:     true,
			Internal:      false},
		{
			Name:          "Marschmusikwertung",
			Date:          time.Date(2018, 9, 17, 0, 0, 0, 0, time.Local),
			Uniform:       "Uniform mit Hut",
			MusicianPlace: "Musikheim",
			Place:         "Spannberg",
			Time:          time.Date(1, 1, 1, 13, 0, 0, 0, time.Local),
			MusicianTime:  time.Date(0, 0, 0, 11, 0, 0, 0, time.Local),
			Important:     false,
			Internal:      false},
		{
			Name:          "Generalversammlung",
			Date:          time.Date(2018, 3, 17, 0, 0, 0, 0, time.Local),
			MusicianPlace: "Musikheim",
			MusicianTime:  time.Date(0, 0, 0, 20, 0, 0, 0, time.Local),
			Important:     true,
			Internal:      true}}

	flg := &model.Instrument{
		NamePlural: "Flügelhörner",
		Name:       "Flügelhorn"}
	hrn := &model.Instrument{
		Name:       "Waldhorn",
		NamePlural: "Waldhörner"}
	ten := &model.Instrument{NamePlural: "Tenorhörner", Name: "Tenorhorn"}

	willi := &model.Member{
		FirstName:    "Willi",
		LastName:     "Herok",
		Joined:       1993,
		Active:       true,
		Deleted:      false,
		LoginAllowed: true,
		Instrument:   flg}

	felix := &model.Member{
		FirstName:    "Felix",
		LastName:     "Nentwich",
		Joined:       1975,
		Active:       true,
		Deleted:      false,
		LoginAllowed: true,
		Instrument:   ten}

	david := &model.Member{
		FirstName:    "David",
		LastName:     "Hörler",
		Joined:       2014,
		Active:       true,
		Deleted:      false,
		LoginAllowed: true,
		Instrument:   flg}

	leo := &model.Member{
		FirstName:    "Leonard",
		LastName:     "Kovarik",
		Joined:       2011,
		Active:       false,
		Deleted:      true,
		LoginAllowed: false,
		Instrument:   hrn}

	obm := &model.LeaderRole{
		Name:       "Obmann",
		NamePlural: "Obmänner"}

	arc := &model.LeaderRole{
		Name:       "Archivar",
		NamePlural: "Archivare"}

	alk := &model.Role{
		Name:       "Alkoholrat",
		NamePlural: "Alkoholräte"}

	willob := &model.LeaderRoleMember{
		LeaderRole: obm,
		Member:     willi,
		Priority:   0}

	felarc := &model.LeaderRoleMember{
		LeaderRole: arc,
		Member:     felix,
		Priority:   0}

	willarc := &model.LeaderRoleMember{
		LeaderRole: arc,
		Member:     willi,
		Priority:   1}

	felalk := &model.RoleMember{
		Role:   alk,
		Member: felix}

	members := []*model.Member{willi, felix, david, leo}
	instruments := []*model.Instrument{flg, hrn, ten}
	roles := []*model.Role{alk}
	leaderRoles := []*model.LeaderRole{obm, arc}
	roleMembers := []*model.RoleMember{felalk}
	leaderRoleMembers := []*model.LeaderRoleMember{willarc, willob, felarc}

	for _, v := range events {
		database.GenericSave(v)
	}
	for _, v := range instruments {
		database.GenericSave(v)
	}
	for _, v := range roles {
		database.GenericSave(v)
	}
	for _, v := range leaderRoles {
		database.GenericSave(v)
	}
	for _, v := range members {
		database.GenericSave(v)
	}

	for _, v := range roleMembers {
		database.GenericSave(v)
	}

	for _, v := range leaderRoleMembers {
		database.GenericSave(v)
	}
}
