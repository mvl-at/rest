package model

type LeaderRole struct {
	ID         int64 `mvlrest:"pk"`
	Name       string
	NamePlural string
}

type LeaderRoleMember struct {
	LeaderRole LeaderRole `mvlrest:"pk"`
	Member     Member     `mvlrest:"pk"`
	Priority   int
}
