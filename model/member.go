package model

import "github.com/mvl-at/qbs"

type Member struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"firstName" roles:"member"`
	LastName     string `json:"lastName" roles:"member"`
	Joined       int    `json:"joined" roles:"member"`
	Picture      string `json:"picture" roles:"member"`
	Active       bool   `json:"active" roles:"member"`
	Deleted      bool   `json:"deleted" roles:"member"`
	LoginAllowed bool   `json:"loginAllowed"`
	Username     string `json:"username"`
	Password     string `json:"-"`

	Instrument   *Instrument `json:"instrument" roles:"member"`
	InstrumentId int64       `qbs:"fk:Instrument" json:"instrumentId" roles:"member"`
}

func (m *Member) Validate(qbs *qbs.Qbs) error {

	if m.Instrument != nil {
		m.InstrumentId = m.Instrument.Id
	}
	return nil
}
