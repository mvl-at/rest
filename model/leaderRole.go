package model

type LeaderRole struct {
	Id         int64
	Name       string
	NamePlural string
}

type LeaderRoleMember struct {
	Id           int64
	LeaderRole   *LeaderRole
	LeaderRoleId int64 `qbs:"fk:LeaderRole"`
	Member       *Member
	MemberId     int64 `qbs:"fk:Member"`
	Priority     int
}
