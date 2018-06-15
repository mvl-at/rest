package model

//Defines an instrument.
type Instrument struct {
	Id         int64  `json:"id"`
	Name       string `json:"name" roles:"instrument"`
	NamePlural string `json:"namePlural" roles:"instrument"`
}
