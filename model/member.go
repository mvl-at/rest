package model

type Member struct {
	Id           int64
	FirstName    string
	LastName     string
	Joined       int
	Picture      string
	Active       bool
	Deleted      bool
	LoginAllowed bool

	Instrument   *Instrument
	InstrumentId int64 `qbs:"fk:Instrument"`
}
