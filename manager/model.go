package manager

import (
	"encoding/json"
	"log"
	"net/http"
	"rest/database"
	"rest/model"
)

type ModelHolder struct {
	Events      []model.Event
	Instruments []model.Instrument
	LeaderRoles []model.LeaderRole
	Members     []model.Member
	Roles       []model.Role
}

var modelHolder ModelHolder

func Init() {
	Routes()
	modelHolder = ModelHolder{
		Events: []model.Event{
			{Name: "baum"},
			{Name: "Frühschoppen", Internal: true}}}
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func Routes() {
	http.HandleFunc("/events", rest(events))
	http.HandleFunc("/members", rest(members))
	http.HandleFunc("/instruments", rest(instruments))
	http.HandleFunc("/roles", rest(roles))
	http.HandleFunc("/events{id}", rest(event))
}

func rest(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("content-type", "application/json")
		next.ServeHTTP(writer, request)
	}
}

func events(rw http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		ev := make([]*model.Event, 0)
		database.GenericFetch(ev)
		err := json.NewEncoder(rw).Encode(ev)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func instruments(rw http.ResponseWriter, r *http.Request) {

	success := false

	if r.Method == "GET" {
		ev := make([]*model.Instrument, 0)
		database.GenericFetch(ev)
		err := json.NewEncoder(rw).Encode(ev)
		if err != nil {
			log.Println(err.Error())
		} else {
			success = true
		}
	}

	if r.Method == "POST" {
		instrument := model.Instrument{}
		err := json.NewDecoder(r.Body).Decode(&instrument)
		if err != nil {
			log.Println(err.Error())
		} else {
			success = true
			database.GenericSave(instrument, true)
		}
	}

	if !success {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func members(rw http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		ev := make([]*model.Member, 0)
		database.GenericFetch(ev)
		err := json.NewEncoder(rw).Encode(ev)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func roles(rw http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		ev := make([]*model.Role, 0)
		database.GenericFetch(ev)
		err := json.NewEncoder(rw).Encode(ev)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func event(rw http.ResponseWriter, r *http.Request) {

}
