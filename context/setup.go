package context

import (
	"log"
	"os"
	"encoding/json"
)

const ConfigPath = "conf.json"

var Log = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var ErrLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
var Conf = config()

func config() (conf *Configuration) {

	conf = &Configuration{}
	fil, err := os.OpenFile(ConfigPath, 0, 0644)
	defer fil.Close()

	if err != nil {
		fil, err = os.Create(ConfigPath)
		defer fil.Close()
		conf = &Configuration{
			Host: "0.0.0.0",
			Port: 8080}
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

type Configuration struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}
