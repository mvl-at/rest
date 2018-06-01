package model

type Instrument struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
}
