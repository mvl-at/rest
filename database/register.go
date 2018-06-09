package database

import (
	"rest/model"
)

func Register() {
	TableCreate(&model.Event{})
	TableCreate(&model.Instrument{})
	TableCreate(&model.Member{})
	TableCreate(&model.LeaderRole{})
	TableCreate(&model.LeaderRoleMember{})
	TableCreate(&model.Role{})
	TableCreate(&model.RoleMember{})
}
