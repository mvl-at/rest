package context

import (
	"os"
	"log"
)

var Log = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var ErrLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)