package api

import (
	"encoding/json"
	"github.com/mvl-at/model"
	"github.com/mvl-at/rest/context"
	"github.com/mvl-at/rest/database"
	"github.com/mvl-at/rest/httpUtils"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const numbers = "0123456789"
const bookTemplate = "book-index.tmp"

//Registers all routes to the http service.
func Handler() http.Handler {
	mux := &http.ServeMux{}
	mux.HandleFunc("/events", httpUtils.Rest(events))
	mux.HandleFunc("/members", httpUtils.Rest(members))
	mux.HandleFunc("/instruments", httpUtils.Rest(instruments))
	mux.HandleFunc("/roles", httpUtils.Rest(roles))
	mux.HandleFunc("/rolesMembers", httpUtils.Rest(rolesMembers))
	mux.HandleFunc("/leaderRoles", httpUtils.Rest(leaderRoles))
	mux.HandleFunc("/leaderRolesMembers", httpUtils.Rest(leaderRolesMembers))
	mux.HandleFunc("/archive", httpUtils.Rest(archive))
	mux.HandleFunc("/login", httpUtils.Rest(login))
	mux.HandleFunc("/credentials", httpUtils.Rest(credentials))
	mux.HandleFunc("/eventsrange", httpUtils.Rest(eventsRange))
	mux.HandleFunc("/userinfo", httpUtils.Rest(userInfo))
	mux.HandleFunc("/bookIndex", httpUtils.Cors(bookIndex))
	return mux
}

//Generic handler for http GET.
//Returns false if the request was not the GET method.
func httpGet(rw http.ResponseWriter, r *http.Request, a interface{}) bool {

	if r.Method == http.MethodGet {
		collection := reflect.New(reflect.SliceOf(reflect.TypeOf(a)))
		database.FindAll(collection.Interface())
		err := json.NewEncoder(rw).Encode(collection.Interface())
		if err != nil {
			context.Log.Println(err.Error())
			return false
		}
		return true
	}
	return false
}

//Generic handler for http POST and PUT.
//Returns false if the requested was neither the POST or PUT method.
func httpPostPut(rw http.ResponseWriter, r *http.Request, a interface{}) (called bool) {

	called = r.Method == http.MethodPost || r.Method == http.MethodPut
	if called {

		token := r.Header.Get("Access-token")
		valid, member, _ := database.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		roles := make([]*model.RoleMember, 0)
		database.FindAllWhereEqual(&roles, "member_id", member.Id)
		modifiedValue := reflect.New(reflect.TypeOf(a).Elem())
		modified := modifiedValue.Interface()
		modifiedRaw, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(modifiedRaw, modified)
		databaseValue := reflect.New(reflect.TypeOf(a).Elem())
		databaseValue.Elem().Set(modifiedValue.Elem())
		databaseEntity := databaseValue.Interface()
		if !database.Find(databaseEntity) {
			databaseEntity = a
			databaseValue = reflect.ValueOf(databaseEntity)
		}
		anyFieldChanges := false
		for i := 0; i < databaseValue.Elem().NumField(); i++ {
			definedRoles := databaseValue.Elem().Type().Field(i).Tag.Get("roles")
			if hasRole(roles, definedRoles) && databaseValue.Elem().Type().Field(i).Tag.Get("json") != "-" {
				databaseValue.Elem().Field(i).Set(modifiedValue.Elem().Field(i))
				anyFieldChanges = true
			}
		}
		if anyFieldChanges {
			database.Save(databaseEntity)
		}
	}
	return
}

//Generic http DELETE method.
//Returns false, if the request was not http DELETE.
func httpDelete(rw http.ResponseWriter, r *http.Request, a interface{}) (called bool) {

	called = r.Method == http.MethodDelete
	if called {

		token := r.Header.Get("access-token")
		valid, member, _ := database.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		roles := make([]*model.RoleMember, 0)
		database.FindAllWhereEqual(&roles, "member_id", member.Id)
		modified := reflect.New(reflect.TypeOf(a).Elem()).Interface()
		modifiedRaw, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(modifiedRaw, modified)
		allowedFields := 0
		modifiedValue := reflect.ValueOf(modified)
		for i := 0; i < modifiedValue.Elem().NumField(); i++ {
			definedRoles := modifiedValue.Elem().Type().Field(i).Tag.Get("roles")
			if definedRoles == "" || hasRole(roles, definedRoles) {
				allowedFields++
			}
		}
		if allowedFields >= modifiedValue.Elem().NumField() {
			database.Delete(modified)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
	}
	return
}

//Checks, if the given member-role association contains at least one of the defined roles.
//Returns true, if either the association has at least one of the defined roles, or if it has the root role.
func hasRole(memberRoles []*model.RoleMember, definedRoles string) bool {

	for _, definedRole := range strings.Split(definedRoles, ",") {
		for _, memberRole := range memberRoles {
			if memberRole.RoleId == definedRole || memberRole.RoleId == "root" {
				return true
			}
		}
	}
	return false
}

//Handler for events.
func events(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Event{}) && !httpPostPut(rw, r, &model.Event{}) && !httpDelete(rw, r, &model.Event{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for instruments.
func instruments(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Instrument{}) && !httpPostPut(rw, r, &model.Instrument{}) && !httpDelete(rw, r, &model.Instrument{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for members.
func members(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Member{}) && !httpPostPut(rw, r, &model.Member{Active: true, LoginAllowed: false}) && !httpDelete(rw, r, &model.Member{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for roles.
func roles(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.Role{}) && !httpPostPut(rw, r, &model.Role{}) && !httpDelete(rw, r, &model.Role{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for leader roles.
func leaderRoles(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.LeaderRole{}) && !httpPostPut(rw, r, &model.LeaderRole{}) && !httpDelete(rw, r, &model.LeaderRole{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for leader roles members.
func leaderRolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.LeaderRoleMember{}) && !httpPostPut(rw, r, &model.LeaderRoleMember{}) && !httpDelete(rw, r, &model.LeaderRoleMember{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for roles members.
func rolesMembers(rw http.ResponseWriter, r *http.Request) {

	if !httpGet(rw, r, &model.RoleMember{}) && !httpPostPut(rw, r, &model.RoleMember{}) && !httpDelete(rw, r, &model.RoleMember{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for archive.
func archive(rw http.ResponseWriter, r *http.Request) {
	ok, member, _ := database.Check(r.Header.Get("Access-Token"))
	if !ok {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	roles := make([]*model.RoleMember, 0)
	database.FindAllWhereEqual(&roles, "member_id", member.Id)
	if !hasRole(roles, "archive") {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	if !httpGet(rw, r, &model.Archive{}) && !httpPostPut(rw, r, &model.Archive{}) && !httpDelete(rw, r, &model.Archive{}) {
		rw.WriteHeader(http.StatusNotFound)
	}
}

//Handler for login.
func login(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	jwtData := database.JWTData{}
	err := json.NewDecoder(r.Body).Decode(&jwtData)

	if err != nil {
		context.Log.Println(err.Error())
	} else {
		success, token := database.Login(&jwtData)

		if success {
			rw.Header().Set("Access-token", token)
			return
		}
	}
	rw.WriteHeader(http.StatusForbidden)
}

//Handler for updating user credentials.
func credentials(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		rw.WriteHeader(http.StatusNotFound)
	}
	token := r.Header.Get("Access-token")
	valid, member, _ := database.Check(token)
	if !valid || member == nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	credentials := &database.Credentials{}
	json.NewDecoder(r.Body).Decode(credentials)
	roles := make([]*model.RoleMember, 0)
	database.FindAllWhereEqual(&roles, "member_id", member.Id)
	if hasRole(roles, "credentials") {
		database.UpdateCredentials(credentials)
	} else {
		if member.Id == credentials.MemberId {
			credentials.Username = member.Username
			database.UpdateCredentials(credentials)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
	}
}

//Handler for events in a certain range
func eventsRange(rw http.ResponseWriter, r *http.Request) {
	fromString := r.URL.Query().Get("from")
	toString := r.URL.Query().Get("to")
	if fromString == "" {
		fromString = "00000101"
	}
	if toString == "" {
		toString = "99991231"
	}
	from, err := time.Parse("20060102", fromString)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	to, err := time.Parse("20060102", toString)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	events := make([]*model.Event, 0)
	database.FindEventsRange(&events, from, to)
	err = json.NewEncoder(rw).Encode(&events)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

//returns information about a user using it's jwt
func userInfo(rw http.ResponseWriter, r *http.Request) {
	jwt := r.Header.Get("Access-token")
	valid, member, _ := database.Check(jwt)
	if !valid {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	userInfo := &database.UserInfo{Member: member, Roles: make([]*model.Role, 0)}
	rolesMembers := make([]*model.RoleMember, 0)
	database.FindAllWhereEqual(&rolesMembers, "member_id", member.Id)
	for _, role := range rolesMembers {
		userInfo.Roles = append(userInfo.Roles, role.Role)
	}
	json.NewEncoder(rw).Encode(&userInfo)
}

// create book indices, only authorized request are able to use a template from a post body
// the book to use can be defined via the comma separated locations http header
func bookIndex(rw http.ResponseWriter, r *http.Request) {
	validMethod := false
	locations := context.Conf.DefaultBooks
	var temp *template.Template = nil
	if r.Method == http.MethodPost {
		token := r.Header.Get("Access-token")
		valid, member, _ := database.Check(token)

		if !valid || member == nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		if len(r.Header.Get("locations")) > 0 {
			locations = strings.Split(r.Header.Get("locations"), ",")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			context.Log.Printf("Cannot read request body: %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		temp, err = template.New("bookIndex").Parse(string(body))
		if err != nil {
			context.Log.Printf("Got an invalid template for a book index: %s\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			_, err := rw.Write([]byte(err.Error()))
			if err != nil {
				context.ErrLog.Println("Cannot write error message to response")
			}
			return
		}
		validMethod = true
	}

	if r.Method == http.MethodGet {
		var err error = nil
		temp, err = template.ParseFiles(bookTemplate)
		if err != nil {
			context.ErrLog.Printf("Cannot parse book template from default file: %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.Header().Set("Content-Disposition", ": inline; filename=\"index."+context.Conf.DefaultBookType+"\"")
		}
		validMethod = true
	}

	if validMethod {
		err := streamLocations(locations, temp, rw)
		if err != nil {
			context.Log.Printf("Cannot generate book index from template: %s\n", err.Error())
			rw.WriteHeader(http.StatusUnprocessableEntity)
		}
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func bookIndexGet()

func streamLocations(locations []string, temp *template.Template, w io.Writer) error {
	index := BookIndex{}
	index.Books = make([]Book, len(locations))
	for i := range index.Books {
		index.Books[i].Scores = make([]*model.Archive, 0)
		index.Books[i].Title = locations[i]
		database.FindAllWhereEqual(&index.Books[i].Scores, "location", locations[i])
		sortScores(index.Books[i].Scores)
	}
	return temp.Execute(w, index)
}

func sortScores(scores []*model.Archive) {
	sort.Slice(scores, func(i, j int) bool {
		// page number is typically stored in the `Note` field
		pageI := strings.Split(scores[i].Note, "-")[0]
		pageJ := strings.Split(scores[j].Note, "-")[0]
		indexI := strings.IndexAny(pageI, numbers)
		indexJ := strings.IndexAny(pageJ, numbers)
		// i has a prefix but j not
		if indexI > 0 && indexJ == 0 {
			return true
		}
		// j has a prefix but i not
		if indexJ > 0 && indexI == 0 {
			return false
		}
		numI, _ := strconv.ParseUint(pageI[indexI:], 10, 64)
		numJ, _ := strconv.ParseUint(pageJ[indexJ:], 10, 64)
		// only comparison of parsable number is required if prefix is the same
		if pageI[:indexI] == pageJ[:indexJ] {
			return numI < numJ
		}

		return strings.Compare(pageI[:indexI], pageJ[:indexJ]) < 0
	})
}
