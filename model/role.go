package model

type Role struct {
	ID         int64 `mvlrest:"pk"`
	Name       string
	NamePlural string
}

type RoleMember struct {
	Role   Role   `mvlrest:"pk"`
	Member Member `mvlrest:"pk"`
}
