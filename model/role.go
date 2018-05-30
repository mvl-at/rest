package model

import "github.com/mvl-at/qbs"

type Role struct {
	Id         int64 `qbs:"pk"`
	Name       string
	NamePlural string
}

type RoleMember struct {
	Role     *Role
	RoleId   int64 `qbs:"fk:Role"`
	Member   *Member
	MemberId int64 `qbs:"fk:Member"`
}

func (*RoleMember) Indexes(indexes *qbs.Indexes) {
	indexes.AddUnique("role_id", "member_id")
}
