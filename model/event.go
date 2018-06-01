package model

import "time"

type Event struct {
	Id            int64     `json:"id"`
	Date          time.Time `json:"date"`
	Time          time.Time `json:"time"`
	MusicianTime  time.Time `json:"musicianTime"`
	Name          string    `json:"name"`
	Note          string    `json:"note"`
	Uniform       string    `json:"uniform"`
	Place         string    `json:"place"`
	MusicianPlace string    `json:"musicianPlace"`
	Internal      bool      `json:"internal"`
	Important     bool      `json:"important"`
}
