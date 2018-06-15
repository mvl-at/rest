package model

import "github.com/mvl-at/qbs"

//Defines a role.
type Role struct {
	Id         string `qbs:"pk" json:"id"`
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
}

//Defines the relation between a role and a member.
type RoleMember struct {
	Id       int64   `json:"id"`
	Role     *Role   `json:"role"`
	RoleId   string  `qbs:"fk:Role" json:"roleId"`
	Member   *Member `json:"member"`
	MemberId int64   `qbs:"fk:Member" json:"memberId"`
}

//Validates all association pointers and assign its id fields to the one of RoleMember.
func (r *RoleMember) Validate(qbs *qbs.Qbs) error {

	if r.Member != nil {
		r.MemberId = r.Member.Id
	}

	if r.Role != nil {
		r.RoleId = r.Role.Id
	}

	return nil
}
