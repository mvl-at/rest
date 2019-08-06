package simple

import (
	json2 "encoding/json"
	"github.com/mvl-at/rest/context"
	"os"
	"time"
)

type persistenceLocation struct {
	data     func() interface{}
	fileName string
}

func PersistenceRunner() {
	os.MkdirAll(context.Conf.PersistenceLocation, 0755)
	persistenceLocations := []persistenceLocation{
		{
			data:     events,
			fileName: context.Conf.EventsFile},
		{data: members,
			fileName: context.Conf.MembersFile}}

	t := time.NewTicker(time.Minute)
	go func() {
		for range t.C {
			for _, location := range persistenceLocations {
				go func(location persistenceLocation) {
					fi, err := os.OpenFile(context.Conf.PersistenceLocation+location.fileName, os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						context.ErrLog.Println(err.Error())
						return
					}
					encoder := json2.NewEncoder(fi)
					err = encoder.Encode(location.data())
					if err != nil {
						context.ErrLog.Println(err.Error())
					}
					if err = fi.Close(); err != nil {
						context.ErrLog.Println(err.Error())
					}
				}(location)
			}
		}
	}()
}
