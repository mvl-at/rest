package model

type Member struct {
	ID           int64 `mvlrest:"pk"`
	FirstName    string
	LastName     string
	Joined       int
	Picture      string
	Active       bool
	Deleted      bool
	LoginAllowed bool

	Instrument Instrument
}
