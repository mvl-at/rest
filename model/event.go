package model

import "time"

type Event struct {
	Id            int64
	Date          time.Time
	Time          time.Time
	MusicianTime  time.Time
	Name          string `json:"name"`
	Note          string
	Uniform       string
	Place         string
	MusicianPlace string
	Internal      bool
	Important     bool
}
