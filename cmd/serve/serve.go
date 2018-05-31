package main

import (
	_ "github.com/mattn/go-sqlite3"
	"rest/database"
	"rest/manager"
	"rest/model"
	"time"
)

func main() {
	//database.Create()
	manager.Register()
	ev := model.Event{
		Name:          "Früschoppen",
		Uniform:       "ohne Hut",
		Note:          "anschl. Abendessen",
		Place:         "rathaus",
		MusicianPlace: "ms",
		Internal:      false,
		Important:     true,
		Date:          time.Date(2018, 6, 17, 0, 0, 0, 0, time.Local),
		Time:          time.Date(1, 1, 0, 9, 30, 0, 0, time.Local)}
	database.GenericSave(&ev, true)
	database.GenericSave(&ev, true)
	//database.GenericFetch(model.Event{})
	manager.Init()
}
