package model

type Role struct {
	Id         int64 `qbs:"pk"`
	Name       string
	NamePlural string
}

type RoleMember struct {
	Id       int64
	Role     *Role
	RoleId   int64 `qbs:"fk:Role"`
	Member   *Member
	MemberId int64 `qbs:"fk:Member"`
}
