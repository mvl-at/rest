package mock

import (
	"rest/database"
	"rest/model"
	"time"
)

var Fruehschoppen = &model.Event{
	Name:          "Frühschoppen",
	Date:          time.Date(2018, 6, 17, 0, 0, 0, 0, time.Local),
	Uniform:       "MVL Polo und Lederhose",
	MusicianPlace: "Musikheim",
	Place:         "Musikheim",
	Time:          time.Date(1, 1, 1, 10, 0, 0, 0, time.Local),
	MusicianTime:  time.Date(0, 0, 0, 8, 45, 0, 0, time.Local),
	Note:          "Anschl. Wegräumen",
	Important:     true,
	Internal:      false}

var Marschmusikwertung = &model.Event{
	Name:          "Marschmusikwertung",
	Date:          time.Date(2018, 9, 17, 0, 0, 0, 0, time.Local),
	Uniform:       "Uniform mit Hut",
	MusicianPlace: "Musikheim",
	Place:         "Spannberg",
	Time:          time.Date(1, 1, 1, 13, 0, 0, 0, time.Local),
	MusicianTime:  time.Date(0, 0, 0, 11, 0, 0, 0, time.Local),
	Important:     false,
	Internal:      false}

var Generalversammlung = &model.Event{
	Name:          "Generalversammlung",
	Date:          time.Date(2018, 3, 17, 0, 0, 0, 0, time.Local),
	MusicianPlace: "Musikheim",
	MusicianTime:  time.Date(0, 0, 0, 20, 0, 0, 0, time.Local),
	Important:     true,
	Internal:      true}

var Fluegelhorn = &model.Instrument{NamePlural: "Flügelhörner", Name: "Flügelhorn"}
var Horn = &model.Instrument{Name: "Waldhorn", NamePlural: "Waldhörner"}
var Tenorhorn = &model.Instrument{NamePlural: "Tenorhörner", Name: "Tenorhorn"}

var Paul = &model.Member{
	FirstName:    "Paul",
	LastName:     "Blatt",
	Joined:       1993,
	Active:       true,
	Deleted:      false,
	LoginAllowed: true,
	Instrument:   Fluegelhorn,
	Username:     "willi",
	Password:     "123456"}

var Helmut = &model.Member{
	FirstName:    "Helmut",
	LastName:     "Gras",
	Joined:       1975,
	Active:       true,
	Deleted:      false,
	LoginAllowed: true,
	Instrument:   Tenorhorn,
	Password:     "dfldfg"}
var Karl = &model.Member{
	FirstName:    "Karl",
	LastName:     "Baum",
	Joined:       2014,
	Active:       true,
	Deleted:      false,
	LoginAllowed: true,
	Instrument:   Fluegelhorn,
	Password:     "dfghhj"}
var Franz = &model.Member{
	FirstName:    "Franz",
	LastName:     "Moos",
	Joined:       2011,
	Active:       false,
	Deleted:      true,
	LoginAllowed: false,
	Instrument:   Horn,
	Password:     "giogftr"}
var Josef = &model.Member{
	FirstName:    "Josef",
	LastName:     "Strauch",
	Joined:       2009,
	Active:       true,
	Deleted:      false,
	LoginAllowed: true,
	Instrument:   Tenorhorn,
	Password:     "dfg"}

var Obmann = &model.LeaderRole{Name: "Obmann", NamePlural: "Obmänner"}
var Archivar = &model.LeaderRole{Name: "Archivar", NamePlural: "Archivare"}
var Root = &model.Role{Id: "root", Name: "Root", NamePlural: "Root"}
var Credentials = &model.Role{Id: "credentials", Name: "Login Verwalter", NamePlural: "Login Verwalter"}
var Instrumente = &model.Role{Id: "instrument", Name: "Instrumenten Manager", NamePlural: "Instrumenten Manager"}
var PaulObmann = &model.LeaderRoleMember{LeaderRole: Obmann, Member: Paul, Priority: 0}
var HelmutArchivar = &model.LeaderRoleMember{LeaderRole: Archivar, Member: Helmut, Priority: 0}
var PaulArchivarStellvertreter = &model.LeaderRoleMember{LeaderRole: Archivar, Member: Paul, Priority: 1}
var JosefRoot = &model.RoleMember{Role: Root, Member: Josef}
var PaulInstrumente = &model.RoleMember{Role: Instrumente, Member: Paul}
var Events = &model.Role{Id: "event", Name: "Termine", NamePlural: "Termine"}
var HelmutEvents = &model.RoleMember{Member: Helmut, Role: Events}
var FranzCredentials = &model.RoleMember{Member: Franz, Role: Credentials}

func MockData() {
	events := []*model.Event{Fruehschoppen, Generalversammlung, Marschmusikwertung}
	members := []*model.Member{Paul, Helmut, Karl, Franz, Josef}
	instruments := []*model.Instrument{Fluegelhorn, Horn, Tenorhorn}
	roles := []*model.Role{Root, Instrumente}
	leaderRoles := []*model.LeaderRole{Obmann, Archivar}
	roleMembers := []*model.RoleMember{JosefRoot, PaulInstrumente, HelmutEvents, FranzCredentials}
	leaderRoleMembers := []*model.LeaderRoleMember{PaulArchivarStellvertreter, PaulObmann, HelmutArchivar}

	for _, v := range events {
		database.Save(v)
	}
	for _, v := range instruments {
		database.Save(v)
	}
	for _, v := range roles {
		database.Save(v)
	}
	for _, v := range leaderRoles {
		database.Save(v)
	}
	for _, v := range members {
		database.Save(v)
	}

	for _, v := range roleMembers {
		database.Save(v)
	}

	for _, v := range leaderRoleMembers {
		database.Save(v)
	}
}
