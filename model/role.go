package model

type Role struct {
	Id         string `qbs:"pk" json:"id"`
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
}

type RoleMember struct {
	Id       int64   `json:"id"`
	Role     *Role   `json:"role"`
	RoleId   string  `qbs:"fk:Role" json:"roleId"`
	Member   *Member `json:"member"`
	MemberId int64   `qbs:"fk:Member" json:"memberId"`
}
