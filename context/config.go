package context

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

//Path for the config file.
const ConfigPath = "conf.json"

var Log = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var ErrLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
var Conf = config()

//Reads the config from file and assigns it to the context.Conf
func config() (conf *Configuration) {
	conf = &Configuration{}
	fil, err := os.OpenFile(ConfigPath, 0, 0644)
	defer fil.Close()

	if err != nil {
		fil, err = os.Create(ConfigPath)
		defer fil.Close()
		rand.Seed(time.Now().UnixNano())
		jwtSecret := make([]byte, 8)
		rand.Read(jwtSecret)
		conf = &Configuration{
			Host:                "0.0.0.0",
			Port:                7301,
			SQLiteFile:          "mvl.sqlite",
			JwtSecret:           fmt.Sprintf("%x", jwtSecret),
			JwtExpiration:       30,
			ApiRoute:            "/api/",
			SimpleRoute:         "/simple/",
			PersistenceLocation: "simple",
			EventsFile:          "events.json",
			MembersFile:         "members.json",
			LeadersFile:         "leaders.json"}
		enc := json.NewEncoder(fil)
		enc.SetIndent("", "  ")
		err = enc.Encode(conf)

	} else {
		err = json.NewDecoder(fil).Decode(conf)
	}

	if err != nil {
		ErrLog.Fatalln(err.Error())
	}
	return
}

//Struct which holds the configuration.
type Configuration struct {
	Host                string `json:"host"`
	Port                uint16 `json:"port"`
	SQLiteFile          string `json:"sqliteFile"`
	JwtSecret           string `json:"jwtSecret"`
	JwtExpiration       int    `json:"jwtExpiration"`
	ApiRoute            string `json:"apiRoute"`
	SimpleRoute         string `json:"simpleRoute"`
	PersistenceLocation string `json:"persistenceLocation"`
	EventsFile          string `json:"eventsFile"`
	MembersFile         string `json:"membersFile"`
	LeadersFile         string `json:"leadersFile"`
}
