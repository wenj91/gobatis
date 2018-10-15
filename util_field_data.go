package gobatis

import (
	"fmt"
	"reflect"
	"time"
)

func interfaceToMap(param interface{}) interface{} {
	val := reflect.ValueOf(param)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	fmt.Println("type name:", typ.Name())

	switch typ.Name(){
	case "bool",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"float32",
		"float64",
		"complex64",
		"complex128",
		"string":

	case "Time":
	case "NullString":
	case "NullInt64":
	case "NullBool":
	case "NullFloat64":
	case "NullTime":
	}

	return nil
}

func dataToFieldVal(data interface{}, tp reflect.Type) interface{} {
	typ := tp.Name()
	switch {
	case typ == "bool" ||
		typ == "int" ||
		typ == "int8" ||
		typ == "int16" ||
		typ == "int32" ||
		typ == "int64" ||
		typ == "uint" ||
		typ == "uint8" ||
		typ == "uint16" ||
		typ == "uint32" ||
		typ == "uint64" ||
		typ == "uintptr" ||
		typ == "float32" ||
		typ == "float64" ||
		typ == "complex64" ||
		typ == "complex128":
		if nil != data {
			return data
		}
	case typ == "string":
		if nil != data {
			return string(data.([]byte))
		}
	case typ == "Time":
		if nil != data {
			tm, _ := time.Parse("2006-01-02 15:04:05", string(data.([]byte)))
			return tm
		}
	case typ == "NullString":
		if nil != data {
			return NullString{string(data.([]byte)), true}
		}
	case typ == "NullInt64":
		if nil != data {
			return NullInt64{data.(int64), true}
		}
	case typ == "NullBool":
		if nil != data {
			return NullBool{data.(bool), true}
		}
	case typ == "NullFloat64":
		if nil != data {
			return NullFloat64{data.(float64), true}
		}
	case typ == "NullTime":
		if nil != data {
			tm, _ := time.Parse("2006-01-02 15:04:05", string(data.([]byte)))
			return NullTime{tm, true}
		}
	}

	return nil
}
