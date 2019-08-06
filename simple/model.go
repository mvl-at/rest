package simple

import (
	"fmt"
	"time"
)

const MemberQuery = "select m.first_name, m.joined, m.last_name, m.picture, i.id, i.name, i.name_plural from member m inner join instrument i on m.instrument_id = i.id where m.active = 1 order by i.name, m.joined, m.last_name, m.first_name"
const EventQuery = "select e.date, e.end, e.important, e.internal, e.musician_place, e.musician_time, e.name, e.note, e.open_end, e.place, e.time from event e where e.date >= date('now') order by e.date, e.musician_time"

var Months = map[string]string{
	"1":  "Jänner",
	"2":  "Feber",
	"3":  "März",
	"4":  "April",
	"5":  "Mai",
	"6":  "Juni",
	"7":  "Juli",
	"8":  "August",
	"9":  "September",
	"10": "Oktober",
	"11": "November",
	"12": "Dezember"}

var Weekdays = []string{"Mo.", "Di.", "Mi.", "Do.", "Fr.", "Sa.", "So."}

type DBScan func(...interface{}) error

type DBO interface {
	Prettyfy()
	IsPretty() bool
	Scan(DBScan, *[]DBO) (*DBO, error)
}

type Member struct {
	Description string `json:"description"`
	joined      int
	lastName    string
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Pretty      bool
}

func (m *Member) Prettyfy() {
	m.Name = fmt.Sprintf("%s %s", m.Name, m.lastName)
	if m.joined != 0 {
		m.Description = fmt.Sprintf("Betritt: %d", m.joined)
	}
	m.Pretty = true
}

func (m *Member) IsPretty() bool {
	return m.Pretty
}

type MemberGroup struct {
	Instrument       string `json:"instrument"`
	InstrumentId     int
	InstrumentPlural string
	Members          []Member `json:"members"`
	Pretty           bool
}

func (m *MemberGroup) Prettyfy() {
	if len(m.Members) >= 2 {
		m.Instrument = m.InstrumentPlural
	}
	for i, member := range m.Members {
		if !member.IsPretty() {
			(&m.Members[i]).Prettyfy()
		}
	}
	m.Pretty = true
}

func (m *MemberGroup) IsPretty() bool {
	return m.Pretty
}

func (*MemberGroup) Scan(scan DBScan, data *[]DBO) (*DBO, error) {
	mg := &MemberGroup{Pretty: false, Members: make([]Member, 1)}
	m := Member{Pretty: false}
	err := scan(&m.Name, &m.joined, &m.lastName, &m.Picture, &mg.InstrumentId, &mg.Instrument, &mg.InstrumentPlural)
	if err != nil {
		return nil, err
	}
	mg.Members[0] = m
	if len(*data) > 0 {
		lastElement := (*data)[len(*data)-1].(*MemberGroup)
		if lastElement.InstrumentId == mg.InstrumentId {
			lastElement.Members = append(lastElement.Members, m)
			(*data)[len(*data)-1] = lastElement
			return nil, nil
		}
	}
	dbo := DBO(mg)
	return &dbo, nil
}

type Event struct {
	Begin             string `json:"begin"`
	beginTime         time.Time
	Day               int `json:"day"`
	date              time.Time
	Ending            string `json:"ending"`
	endingTime        time.Time
	Important         bool   `json:"important"`
	Internal          bool   `json:"internal"`
	MusicianBegin     string `json:"musician_begin"`
	musicianBeginTime time.Time
	musicianPlace     string
	Name              string `json:"name"`
	Note              string `json:"note"`
	openEnd           int
	place             string
	Weekday           string `json:"weekday"`
	pretty            bool
}

func (e *Event) Prettyfy() {
	e.Day = e.date.Day()
	e.Begin = fmt.Sprintf("%s Uhr, %s", e.beginTime.Format("15:04"), e.Begin)
	e.MusicianBegin = fmt.Sprintf("%s Uhr, %s", e.musicianBeginTime.Format("15:04"), e.musicianPlace)
	switch e.openEnd {
	case 0:
		e.Ending = fmt.Sprintf("bis ca. %s Uhr", e.endingTime.Format("15:04"))
	case 1:
		e.Ending = "offenes Ende"
	case 2:
		e.Ending = "Ende unbekannt, %s"
	}
	if e.Note != "" {
		e.Ending += ", " + e.Note
	}
	e.Weekday = Weekdays[(int(e.date.Weekday())+6)%7]
	e.pretty = true
}

func (e *Event) IsPretty() bool {
	return e.pretty
}

func (e *Event) Scan(scan DBScan, data *[]DBO) (*DBO, error) {
	return nil, nil
}

type EventGroup struct {
	Events []Event `json:"events"`
	Month  string  `json:"month"`
	Pretty bool
}

func (e *EventGroup) Prettyfy() {
	e.Month = Months[e.Month]
	for i, event := range e.Events {
		if !event.IsPretty() {
			(&e.Events[i]).Prettyfy()
		}
	}
	e.Pretty = true
}

func (e *EventGroup) IsPretty() bool {
	return e.Pretty
}

func (*EventGroup) Scan(scan DBScan, data *[]DBO) (*DBO, error) {
	eg := &EventGroup{Events: make([]Event, 1), Pretty: false}
	e := Event{pretty: false}
	err := scan(&e.date, &e.endingTime, &e.Important, &e.Internal, &e.musicianPlace, &e.musicianBeginTime, &e.Name, &e.Note, &e.openEnd, &e.place, &e.beginTime)
	if err != nil {
		return nil, err
	}
	eg.Events[0] = e
	eg.Month = e.date.Format("1")
	if len(*data) > 0 {
		lastElement := (*data)[len(*data)-1].(*EventGroup)
		if lastElement.Month == eg.Month {
			lastElement.Events = append(lastElement.Events, e)
			(*data)[len(*data)-1] = lastElement
			return nil, nil
		}
	}
	dbo := DBO(eg)
	return &dbo, nil
}
