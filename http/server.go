package http

import (
	"fmt"
	"log"
	"net/http"
	"rest/context"
)

//Serves the http application and logs it.
func Run() {
	host := fmt.Sprintf("%s:%d", context.Conf.Host, context.Conf.Port)
	context.Log.Println("Listen on " + host)
	err := http.ListenAndServe(host, nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}
