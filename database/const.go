package database

import (
	"reflect"
	"time"
)

/**
these vars should be taken as constants as it is not possible to declare not native constants in go
*/

var datatype = map[string]string{
	reflect.String.String():           "text",
	reflect.TypeOf(time.Now()).Name(): "integer",
	reflect.Bool.String():             "boolean",
	reflect.Int.String():              "integer",
	reflect.Int64.String():            "integer"}
