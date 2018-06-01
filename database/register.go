package database

import (
	"rest/model"
)

func Register() {
	GenericCreate(&model.Event{})
	GenericCreate(&model.Instrument{})
	GenericCreate(&model.Member{})
	GenericCreate(&model.LeaderRole{})
	GenericCreate(&model.LeaderRoleMember{})
	GenericCreate(&model.Role{})
	GenericCreate(&model.RoleMember{})
}
