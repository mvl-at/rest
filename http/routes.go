package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"rest/context"
	"rest/database"
	"rest/model"
	"rest/security"
	"strings"
)

func Routes() {
	http.HandleFunc("/events", rest(events))
	http.HandleFunc("/members", rest(members))
	http.HandleFunc("/instruments", rest(instruments))
	http.HandleFunc("/roles", rest(roles))
	http.HandleFunc("/rolesMembers", rest(rolesMembers))
	http.HandleFunc("/leaderRoles", rest(leaderRoles))
	http.HandleFunc("/leaderRolesMembers", rest(leaderRolesMembers))
	http.HandleFunc("/login", rest(login))
}

func rest(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("content-type", "application/json")
		next.ServeHTTP(writer, request)
	}
}

func get(rw http.ResponseWriter, r *http.Request, a interface{}) bool {

	if r.Method == "GET" {
		collection := reflect.New(reflect.SliceOf(reflect.TypeOf(a)))
		database.GenericFetch(collection.Interface())
		err := json.NewEncoder(rw).Encode(collection.Interface())
		if err != nil {
			context.Log.Println(err.Error())
			return false
		}
		return true
	}
	return false
}

func set(rw http.ResponseWriter, r *http.Request, a interface{}) (called bool) {

	called = r.Method == "POST"
	if called {

		token := r.Header.Get("token")
		valid, member := security.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		roles := make([]*model.RoleMember, 0)
		database.GenericFetchWhereEqual(&roles, "member_id", member.Id)
		defaultValue := reflect.ValueOf(a)
		entityValue := reflect.New(reflect.TypeOf(a).Elem())
		entity := entityValue.Elem().Interface()
		rawEntity, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(rawEntity, &entity)

		if err != nil {
			context.ErrLog.Fatalln(err)
		}

		databaseValue := reflect.New(reflect.TypeOf(a).Elem())
		//databaseValue.Elem().Set(reflect.ValueOf(entity).Elem())
		databaseEntity := reflect.ValueOf(entity).Elem().Interface()
		if database.GenericSingleFetch(databaseEntity) {
			defaultValue.Elem().Elem().Set(databaseValue.Elem().Elem())
		}

		for i := 0; i < entityValue.Elem().NumField(); i++ {
			definedRoles := entityValue.Elem().Type().Field(i).Tag.Get("roles")

			if hasRole(roles, definedRoles) {
				defaultValue.Elem().Field(i).Set(entityValue.Elem().Field(i))
			}

		}

		database.GenericSave(defaultValue)
	}
	return
}

func hasRole(memberRoles []*model.RoleMember, definedRoles string) bool {

	for _, definedRole := range strings.Split(definedRoles, ",") {
		for _, memberRole := range memberRoles {
			if memberRole.RoleId == definedRole {
				return true
			}
		}
	}
	return false
}

func events(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.Event{}) && !set(rw, r, &model.Event{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func instruments(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.Instrument{}) && !set(rw, r, &model.Instrument{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func members(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.Member{}) && !set(rw, r, &model.Member{Deleted: false, Active: true, LoginAllowed: false}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func roles(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.Role{}) && !set(rw, r, &model.Role{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func leaderRoles(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.LeaderRole{}) && !set(rw, r, &model.LeaderRole{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func leaderRolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.LeaderRoleMember{}) && !set(rw, r, &model.LeaderRoleMember{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func rolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.RoleMember{}) && !set(rw, r, &model.RoleMember{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func login(rw http.ResponseWriter, r *http.Request) {
	jwtData := security.JWTData{}
	err := json.NewDecoder(r.Body).Decode(&jwtData)

	if err != nil {
		context.Log.Println(err.Error())
	} else {
		success, token := security.Login(&jwtData)

		if success {
			rw.Write([]byte(token))
			return
		}
	}
	rw.WriteHeader(http.StatusForbidden)
}
