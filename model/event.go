package model

import "time"

type Event struct {
	ID            int64 `mvlrest:"pk"`
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
