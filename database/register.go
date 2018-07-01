package database

import (
	"fmt"
	"github.com/mvl-at/model"
	"github.com/mvl-at/rest/context"
	"math/rand"
)

//Registers all model structs from this project to the database.
func Register() {
	TableCreate(&model.Event{})
	TableCreate(&model.Instrument{})
	TableCreate(&model.Member{})
	TableCreate(&model.LeaderRole{})
	TableCreate(&model.LeaderRoleMember{})
	TableCreate(&model.Role{})
	TableCreate(&model.RoleMember{})
}

//Checks if root user and role exists.
//If not, it creates both, the password will be printed on configured log tool.
func CheckRoot() {
	memberRoles := make([]*model.RoleMember, 0)
	FindAllWhereEqual(&memberRoles, "role_id", "root")
	if len(memberRoles) <= 0 {
		password := make([]byte, 8)
		rand.Read(password)

		roles := make([]*model.Role, 0)
		FindAllWhereEqual(&roles, "id", "root")
		var rootRole *model.Role
		if len(roles) <= 0 {
			rootRole = &model.Role{Id: "root", Name: "Root", NamePlural: "Root"}
			Save(rootRole)
			context.Log.Println("role root was not found, so it was created")
		} else {
			rootRole = roles[0]
		}
		root := &model.Member{Username: "root", LoginAllowed: true, Deleted: false, Password: fmt.Sprintf("%x", password)}
		rootRootRole := &model.RoleMember{Member: root, Role: rootRole}
		Save(root)
		Save(rootRootRole)
		context.Log.Printf("no root user was found, so one was automatically created with username %s and password %x. Use it to log in, create another root user and delete this user", root.Username, password)
	}
}
