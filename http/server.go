package http

import (
	"fmt"
	"github.com/mvl-at/rest/api"
	"github.com/mvl-at/rest/context"
	"github.com/mvl-at/rest/simple"
	"log"
	"net/http"
)

//Serves the http application and logs it.
func Run() {
	host := fmt.Sprintf("%s:%d", context.Conf.Host, context.Conf.Port)
	mux := &http.ServeMux{}
	handle(mux, api.Handler(), context.Conf.ApiRoute)
	handle(mux, simple.Handler(), context.Conf.SimpleRoute)
	context.Log.Println("Listen on " + host)
	err := http.ListenAndServe(host, mux)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func handle(mux *http.ServeMux, childHandler http.Handler, url string) {
	mux.Handle(url, http.StripPrefix(url[:len(url)-1], childHandler))
}
