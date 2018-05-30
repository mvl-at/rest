package model

type Instrument struct {
	ID         int64 `mvlrest:"pk"`
	Name       string
	NamePlural string
}
