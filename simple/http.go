package simple

import (
	json2 "encoding/json"
	"fmt"
	"github.com/mvl-at/rest/httpUtils"
	"net/http"
)

func Handler() http.Handler {
	mux := &http.ServeMux{}
	mux.HandleFunc("/events", httpUtils.Rest(json(events)))
	mux.HandleFunc("/members", httpUtils.Rest(json(members)))
	return mux
}

func json(data func() interface{}) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		encoder := json2.NewEncoder(writer)
		err := encoder.Encode(data())
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func events() interface{} {
	events := []DBO{&EventGroup{}}
	err := QueryData(EventQuery, &events, 11)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range events {
		v.Prettyfy()
	}
	return events
}

func members() interface{} {
	members := []DBO{&MemberGroup{}}
	err := QueryData(MemberQuery, &members, 7)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range members {
		v.Prettyfy()
	}
	return members
}
