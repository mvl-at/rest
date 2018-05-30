package model

import "github.com/mvl-at/qbs"

type LeaderRole struct {
	Id         int64
	Name       string
	NamePlural string
}

type LeaderRoleMember struct {
	LeaderRole   *LeaderRole
	LeaderRoleId int64 `qbs:"fk:LeaderRole"`
	Member       *Member
	MemberId     int64 `qbs:"fk:Member"`
	Priority     int
}

func (*LeaderRoleMember) Indexes(indexes *qbs.Indexes) {
	indexes.AddUnique("leader_role_id", "member_id")
}
