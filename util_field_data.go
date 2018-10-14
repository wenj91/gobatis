package gobatis

import (
	"fmt"
	"reflect"
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