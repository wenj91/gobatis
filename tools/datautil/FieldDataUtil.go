package datautil

import (
	"reflect"
	"../../structs"
	"time"
	"log"
)

func DataToFieldVal(data interface{}, tp reflect.Type) interface{} {

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
			return structs.NullString{string(data.([]byte)), true}
		}
	case typ == "NullInt64":
		if nil != data {
			return structs.NullInt64{data.(int64), true}
		}
	case typ == "NullBool":
		if nil != data {
			return structs.NullBool{data.(bool), true}
		}
	case typ == "NullFloat64":
		if nil != data {
			return structs.NullFloat64{data.(float64), true}
		}
	case typ == "NullTime":
		if nil != data {
			tm, _ := time.Parse("2006-01-02 15:04:05", string(data.([]byte)))
			return structs.NullTime{tm, true}
		}
	}

	return nil
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func FieldToParams(param interface{}, fieldName string) interface{} {
	paramVal := reflect.ValueOf(&param)
	if paramVal.Kind() != reflect.Ptr {
		log.Fatal("params parse exception")
		return nil
	}

	ptr := reflect.Indirect(paramVal).Elem()
	fieldVal := ptr.FieldByName(fieldName)

	if IsZeroOfUnderlyingType(fieldVal) {
		log.Fatal("no this field")
		return nil
	}

	typ := fieldVal.Type().Name()
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
		typ == "complex128" ||
		typ == "string":
		return fieldVal.Interface()
	case typ == "Time":
		tm := fieldVal.Interface().(time.Time)
		return tm.Format("2006-01-02 15:04:05")
	case typ == "NullString":
		obj := fieldVal.Interface().(structs.NullString)
		if obj.Valid {
			return obj.String
		}
	case typ == "NullInt64":
		obj := fieldVal.Interface().(structs.NullInt64)
		if obj.Valid {
			return obj.Int64
		}
	case typ == "NullBool":
		obj := fieldVal.Interface().(structs.NullBool)
		if obj.Valid {
			return obj.Bool
		}
	case typ == "NullFloat64":
		obj := fieldVal.Interface().(structs.NullFloat64)
		if obj.Valid {
			return obj.Float64
		}
	case typ == "NullTime":
		obj := fieldVal.Interface().(structs.NullTime)
		if obj.Valid {
			return obj.Time.Format("2006-01-02 15:04:05")
		}
	}

	return nil
}