package simple

type Member struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
}

type MemberGroup struct {
	Instrument string   `json:"instrument"`
	Members    []Member `json:"members"`
}

type Event struct {
	Begin         string `json:"begin"`
	Day           int    `json:"day"`
	Ending        string `json:"ending"`
	MusicianBegin string `json:"musician_begin"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	Weekday       string `json:"weekday"`
}

type EventGroup struct {
	Events []Event `json:"events"`
	Month  string  `json:"month"`
}
