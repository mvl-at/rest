package http

import (
	"encoding/json"
	"log"
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

func set(rw http.ResponseWriter, r *http.Request, a interface{}) bool {

	if r.Method == "POST" {

		token := r.Header.Get("token")
		valid, member := security.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return false
		}
		roles := make([]*model.RoleMember, 0)
		database.GenericFetchWhereEqual(roles, "member_id", member.Id)
		defaultValue := reflect.ValueOf(a)
		entity := reflect.New(reflect.TypeOf(a))
		rawEntity := make([]byte, 0)
		r.Body.Read(rawEntity)
		json.Unmarshal(rawEntity, entity)

		databaseEntity := reflect.New(reflect.TypeOf(a))
		if database.GenericSingleFetch(databaseEntity) {
			defaultValue = databaseEntity
		}

		for i := 0; i < entity.Elem().NumField(); i++ {
			definedRoles := entity.Elem().Type().Field(i).Tag.Get("roles")

			if hasRole(roles, definedRoles) {
				defaultValue.Elem().Field(i).Set(entity.Elem().Field(i))
			}

		}

		database.GenericSave(defaultValue)
	}
	return true
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

	if !get(rw, r, &model.Event{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func instruments(rw http.ResponseWriter, r *http.Request) {

	success := get(rw, r, &model.Instrument{})

	if r.Method == "POST" {
		instrument := model.Instrument{}
		err := json.NewDecoder(r.Body).Decode(&instrument)
		if err != nil {
			log.Println(err.Error())
		} else {
			success = true
			database.GenericSave(instrument)
		}
	}

	if !success {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func members(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.Member{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func roles(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.Role{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func leaderRoles(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.LeaderRole{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func leaderRolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.LeaderRoleMember{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func rolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !get(rw, r, &model.RoleMember{}) {
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
