package model

type Member struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Joined       int    `json:"joined"`
	Picture      string `json:"picture"`
	Active       bool   `json:"active"`
	Deleted      bool   `json:"deleted"`
	LoginAllowed bool   `json:"loginAllowed"`

	Instrument   *Instrument `json:"instrument"`
	InstrumentId int64       `qbs:"fk:Instrument"json:"instrumentId"`
}
