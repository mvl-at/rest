package model

import "github.com/mvl-at/qbs"

//Defines a leader role.
type LeaderRole struct {
	Id         int64  `json:"id"`
	Name       string `json:"name" roles:"leader"`
	NamePlural string `json:"namePlural" roles:"leader"`
}

//Defines the relation between a leader role an a member.
type LeaderRoleMember struct {
	Id           int64       `json:"id"`
	LeaderRole   *LeaderRole `json:"leaderRole" roles:"leader"`
	LeaderRoleId int64       `qbs:"fk:LeaderRole" json:"leaderRoleId" roles:"leader"`
	Member       *Member     `json:"member" roles:"leader"`
	MemberId     int64       `qbs:"fk:Member" json:"memberId" roles:"leader"`
	Priority     int         `json:"priority" roles:"leader"`
}

//Validates all association pointers and assign its id fields to the one of LeaderRoleMember.
func (l *LeaderRoleMember) Validate(qbs *qbs.Qbs) error {

	if l.Member != nil {
		l.MemberId = l.Member.Id
	}

	if l.LeaderRole != nil {
		l.LeaderRoleId = l.LeaderRole.Id
	}

	return nil
}
