package simple

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "mvl.sqlite")
	if err != nil {
		return err
	}
	return err
}

//func Events() []Event {
//	db.Query("select * from event")
//}

func Leaders() ([]Member, error) {
	rows, err :=
}

func Members() ([]MemberGroup, error) {
	rows, err := db.Query("select m.first_name, m.last_name, m.joined, m.picture, i.id, i.name, i.name_plural from member m inner join instrument i on m.instrument_id = i.id order by i.name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	members := make([]MemberGroup, 0)
	var firstName string
	var lastName string
	var joined int
	var picture string
	var instrumentId int
	var instrumentName string
	var instrumentNamePlural string
	oldInstrumentId := -1
	currentIndex := -1
	member := &Member{}
	for rows.Next() {
		rows.Scan(&firstName, &lastName, &joined, &picture, &instrumentId, &instrumentName, &instrumentNamePlural)
		member.Name = fmt.Sprintf("%s %s", firstName, lastName)
		member.Description = fmt.Sprintf("Beitritt: %d", joined)
		member.Picture = picture
		if instrumentId != oldInstrumentId {
			currentIndex++
			memberGroup := &MemberGroup{Members: make([]Member, 1)}
			memberGroup.Instrument = instrumentName
			memberGroup.Members[0] = *member
			members = append(members, *memberGroup)
		} else {
			members[currentIndex].Members = append(members[currentIndex].Members, *member)
			members[currentIndex].Instrument = instrumentNamePlural
		}
		oldInstrumentId = instrumentId
	}
	return members, nil
}
