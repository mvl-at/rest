package database

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

func GenericCreate(a interface{}) {
	t := reflect.TypeOf(a)

	tableName := t.Name() + "s"
	fieldList := ""

	keys := 0
	keyList := ",primary key("

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("mvlrest") == "pk" {
			keys++
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldList += field.Name + " " + datatype[field.Type.String()]

		if field.Tag.Get("mvlrest") == "pk" {

			if keys == 1 {
				fieldList += " primary key"
			} else {
				keyList += field.Name + ","
			}
		}

		if i != t.NumField()-1 {
			fieldList += ","
		}
	}

	if keys > 1 {
		keyList = keyList[:len(keyList)-1] + ")"
	} else {
		keyList = ""
	}

	sqlCmd := strings.ToLower(fmt.Sprintf("create table if not exists %s(%s %s)", tableName, fieldList, keyList))
	log.Println(sqlCmd)
	_, err := db.Exec(sqlCmd)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func GenericDelete(a interface{}) {
	v := reflect.ValueOf(a).Elem()
	tableName := strings.ToLower(v.Type().Name() + "s")
	fieldMap := columnPairs(v, true, false)
	fields, values := partMap(fieldMap, false)
	fieldsString := chain(fields, values, "= ?", "and", true, false)

	sqlCmd := strings.ToLower(fmt.Sprintf("delete from %s where%s", tableName, fieldsString))
	log.Println(sqlCmd)
	_, err := db.Exec(sqlCmd, values...)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GenericSave(a interface{}, autoKey bool) {
	if GenericExists(a) {
		update(a)
	} else {
		insert(a, autoKey)
	}
}

func insert(a interface{}, autoKey bool) {
	v := reflect.ValueOf(a)
	tableName := tableName(v)
	fieldMap := columnPairs(v, !autoKey, true)
	fields, values := partMap(fieldMap, false)
	fieldsString := chain(fields, values, "", ",", true, false)

	sqlCmd := strings.ToLower(fmt.Sprintf("insert into %s (%s) values (%s)", tableName, fieldsString, argChain(len(values), "?", ",")))
	log.Println(sqlCmd)
	_, err := db.Exec(sqlCmd, values...)
	if err != nil {
		log.Fatal(err.Error())
	}
}
func update(a interface{}) {
	v := reflect.ValueOf(a)
	tableName := tableName(v)
	fieldMap := columnPairs(v, false, true)
	fields, values := partMap(fieldMap, false)
	fieldsString := chain(fields, values, "= ?", ",", true, false)

	sqlCmd := strings.ToLower(fmt.Sprintf("update %s set%s", tableName, fieldsString))
	log.Println(sqlCmd)
	_, err := db.Exec(sqlCmd, values...)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func serialize(v reflect.Value) interface{} {

	if v.Type().Kind() == reflect.TypeOf(time.Now()).Kind() {
		return v.Interface().(time.Time).Unix()
	}
	return v.Interface()
}

func GenericExists(a interface{}) bool {
	v := reflect.ValueOf(a)
	//v=v.Elem()
	tableName := tableName(v)
	fields := columnPairs(v, true, false)
	_, params := partMap(fields, true)
	sqlCmd := strings.ToLower(fmt.Sprintf("select * from %s where%s", tableName, argChain(len(fields), "? = ?", " and")))
	var m string
	err := db.QueryRow(sqlCmd, params...).Scan(&m)

	log.Println(sqlCmd)

	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err.Error())
	}

	return err == nil
}

func GenericFetch(a interface{}) []interface{} {

	ret := make([]interface{}, 0)

	v := reflect.ValueOf(a)
	//v=v.Elem()
	tableName := tableName(v)
	fieldMap := columnPairs(v, true, true)
	fields, _ := partMap(fieldMap, false)

	sqlCmd := strings.ToLower(fmt.Sprintf("select %s from %s", chain(fields, nil, "", ",", true, false), tableName))
	log.Println(sqlCmd)
	rows, err := db.Query(sqlCmd)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		pointers := make([]interface{}, len(fields))
		values := make([]interface{}, len(fields))
		for e := range pointers {
			pointers[e] = &values[e]
		}

		rows.Scan(pointers...)
		for i, v := range values {

			strData, success := v.([]byte)
			if success {
				values[i] = string(strData)
			}
		}
		v := reflect.New(reflect.TypeOf(a))

		for e := range fields {
			fieldName := fields[e]
			value := values[e]
			field := v.Elem().FieldByName(fieldName)
			t := field.Type()
			deserializeValue := deserialize(t, value)

			field.Set(deserializeValue)
		}
		ret = append(ret, v.Elem().Interface())
	}
	return ret
}

func deserialize(kind reflect.Type, a interface{}) reflect.Value {

	raw := a

	if kind == reflect.TypeOf(time.Now()) {
		raw = time.Unix(a.(int64), 0)
	}

	return reflect.ValueOf(raw)
}

func param(old []string) (ret []interface{}) {
	ret = make([]interface{}, len(old))
	for k, v := range old {
		ret[k] = v
	}
	return
}

func tableName(v reflect.Value) string {
	return strings.ToLower(v.Type().Name()) + "s"
}

func argChain(amount int, expression string, delimiter string) (ret string) {
	ret = ""
	for i := 0; i < amount; i++ {
		ret += " " + expression + delimiter
	}
	ret = ret[:len(ret)-len(delimiter)]
	return
}

func chain(keys []string, values []interface{}, operator string, chainOperator string, useFields bool, useValues bool) (ret string) {
	ret = ""
	for i := 0; i < len(keys); i++ {

		field := ""
		var value interface{} = ""

		if useFields {
			field = keys[i]
		}

		if useValues {
			value = values[i]
		}

		ret += fmt.Sprintf(" %s %s %s %s", field, operator, value, chainOperator)
	}
	ret = ret[:len(ret)-len(chainOperator)-1]
	return
}

func columnPairs(v reflect.Value, keys bool, nonKeys bool) (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)
		if (field.Tag.Get("mvlrest") == "pk" && keys) || (nonKeys && field.Tag.Get("mvlrest") != "pk") {
			ret[field.Name] = serialize(fieldValue)
		}
	}
	return
}

func partMap(m map[string]interface{}, combination bool) (keys []string, values []interface{}) {
	keys = make([]string, 0)
	values = make([]interface{}, 0)
	for k, v := range m {

		if combination {
			values = append(values, k)
		} else {
			keys = append(keys, k)
		}
		values = append(values, v)
	}
	return
}
