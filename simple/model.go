package simple

import "fmt"

const MemberQuery = "select m.first_name, m.joined, m.last_name, m.picture, i.id, i.name, i.name_plural from member m inner join instrument i on m.instrument_id = i.id where m.active = 1 order by i.name, m.joined, m.last_name, m.first_name"

type DBScan func(...interface{}) error

type DBO interface {
	Prettyfy()
	IsPretty() bool
	Scan(DBScan, *[]DBO) (*DBO, error)
}

type Member struct {
	Description string `json:"description"`
	Joined      int
	LastName    string
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Pretty      bool
}

func (m *Member) Prettyfy() {
	m.Name = fmt.Sprintf("%s %s", m.Name, m.LastName)
	if m.Joined != 0 {
		m.Description = fmt.Sprintf("Betritt: %d", m.Joined)
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
	err := scan(&m.Name, &m.Joined, &m.LastName, &m.Picture, &mg.InstrumentId, &mg.Instrument, &mg.InstrumentPlural)
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
	Begin         string `json:"begin"`
	Day           int    `json:"day"`
	Ending        string `json:"ending"`
	MusicianBegin string `json:"musician_begin"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	Weekday       string `json:"weekday"`
	Pretty        bool
}

func (e *Event) Prettyfy() {
	panic("implement me")
}

func (e *Event) IsPretty() bool {
	return e.Pretty
}

func (e *Event) Scan(scan DBScan, data *[]DBO) (*DBO, error) {}

type EventGroup struct {
	Events []Event `json:"events"`
	Month  string  `json:"month"`
	Pretty bool
}

func (e *EventGroup) Prettyfy() {
	panic("implement me")
}

func (e *EventGroup) IsPretty() bool {
	return e.Pretty
}

func (e *EventGroup) Scan(scan DBScan, data *[]DBO) *DBO {
	panic("implement me")
}
