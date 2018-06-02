package model

type LeaderRole struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
}

type LeaderRoleMember struct {
	Id           int64       `json:"id"`
	LeaderRole   *LeaderRole `json:"leaderRole"`
	LeaderRoleId int64       `qbs:"fk:LeaderRole" json:"leaderRoleId"`
	Member       *Member     `json:"member"`
	MemberId     int64       `qbs:"fk:Member" json:"memberId"`
	Priority     int         `json:"priority"`
}
