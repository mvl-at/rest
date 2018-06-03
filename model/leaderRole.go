package model

type LeaderRole struct {
	Id         int64  `json:"id"`
	Name       string `json:"name" roles:"leader"`
	NamePlural string `json:"namePlural" roles:"leader"`
}

type LeaderRoleMember struct {
	Id           int64       `json:"id"`
	LeaderRole   *LeaderRole `json:"leaderRole" roles:"leader"`
	LeaderRoleId int64       `qbs:"fk:LeaderRole" json:"leaderRoleId" roles:"leader"`
	Member       *Member     `json:"member" roles:"leader"`
	MemberId     int64       `qbs:"fk:Member" json:"memberId" roles:"leader"`
	Priority     int         `json:"priority" roles:"leader"`
}
