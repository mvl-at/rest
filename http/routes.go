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

func httpGet(rw http.ResponseWriter, r *http.Request, a interface{}) bool {

	if r.Method == http.MethodGet {
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

func httpPostPut(rw http.ResponseWriter, r *http.Request, a interface{}) (called bool) {

	called = r.Method == http.MethodPost || r.Method == http.MethodPut
	if called {

		token := r.Header.Get("token")
		valid, member := security.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		roles := make([]*model.RoleMember, 0)
		database.GenericFetchWhereEqual(&roles, "member_id", member.Id)
		modifiedValue := reflect.New(reflect.TypeOf(a).Elem())
		modified := modifiedValue.Interface()
		modifiedRaw, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(modifiedRaw, modified)
		databaseValue := reflect.New(reflect.TypeOf(a).Elem())
		databaseValue.Elem().Set(modifiedValue.Elem())
		databaseEntity := databaseValue.Interface()
		if !database.GenericSingleFetch(databaseEntity) {
			databaseEntity = a
			databaseValue = reflect.ValueOf(databaseEntity)
		}
		anyFieldChanges := false
		for i := 0; i < databaseValue.Elem().NumField(); i++ {
			definedRoles := databaseValue.Elem().Type().Field(i).Tag.Get("roles")
			if hasRole(roles, definedRoles) {
				databaseValue.Elem().Field(i).Set(modifiedValue.Elem().Field(i))
				anyFieldChanges = true
			}
		}
		if anyFieldChanges {
			database.GenericSave(databaseEntity)
		}
	}
	return
}

func httpDelete(rw http.ResponseWriter, r *http.Request, a interface{}) (called bool) {

	called = r.Method == http.MethodDelete
	if called {

		token := r.Header.Get("token")
		valid, member := security.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		roles := make([]*model.RoleMember, 0)
		database.GenericFetchWhereEqual(&roles, "member_id", member.Id)
		modified := reflect.New(reflect.TypeOf(a).Elem()).Interface()
		modifiedRaw, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(modifiedRaw, modified)
		allowedFields := 0
		modifiedValue := reflect.ValueOf(modified)
		for i := 0; i < modifiedValue.Elem().NumField(); i++ {
			definedRoles := modifiedValue.Elem().Type().Field(i).Tag.Get("roles")
			if hasRole(roles, definedRoles) {
				allowedFields++
			}
		}
		if allowedFields >= modifiedValue.Elem().NumField()-1 {
			database.GenericDelete(modified)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
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

	if !httpGet(rw, r, &model.Event{}) && !httpPostPut(rw, r, &model.Event{}) && !httpDelete(rw, r, &model.Event{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func instruments(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Instrument{}) && !httpPostPut(rw, r, &model.Instrument{}) && !httpDelete(rw, r, &model.Instrument{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func members(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Member{}) && !httpPostPut(rw, r, &model.Member{Deleted: false, Active: true, LoginAllowed: false}) && !httpDelete(rw, r, &model.Member{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func roles(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Role{}) && !httpPostPut(rw, r, &model.Role{}) && !httpDelete(rw, r, &model.Role{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func leaderRoles(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.LeaderRole{}) && !httpPostPut(rw, r, &model.LeaderRole{}) && !httpDelete(rw, r, &model.LeaderRole{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func leaderRolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.LeaderRoleMember{}) && !httpPostPut(rw, r, &model.LeaderRoleMember{}) && !httpDelete(rw, r, &model.LeaderRoleMember{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func rolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.RoleMember{}) && !httpPostPut(rw, r, &model.RoleMember{}) && !httpDelete(rw, r, &model.RoleMember{}) {
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
