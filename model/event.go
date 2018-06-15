package model

import "time"

//Defines an event.
type Event struct {
	Id            int64     `json:"id"`
	Date          time.Time `json:"date" roles:"event"`
	Time          time.Time `json:"time" roles:"event"`
	MusicianTime  time.Time `json:"musicianTime" roles:"event"`
	Name          string    `json:"name" roles:"event"`
	Note          string    `json:"note" roles:"event"`
	Uniform       string    `json:"uniform" roles:"event"`
	Place         string    `json:"place" roles:"event"`
	MusicianPlace string    `json:"musicianPlace" roles:"event"`
	Internal      bool      `json:"internal" roles:"event"`
	Important     bool      `json:"important" roles:"event"`
}
