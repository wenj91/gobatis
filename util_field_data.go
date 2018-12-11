package gobatis

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
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

func bytesToVal(data interface{}, tp reflect.Type) interface{} {
	str := string(data.([]uint8))
	switch tp.Kind() {
	case reflect.Bool:
		if str == "1" {
			data = true
		}else{
			data = false
		}
	case reflect.Int:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int(i)
	case reflect.Int8:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int8(i)
	case reflect.Int16:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int16(i)
	case reflect.Int32:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int32(i)
	case reflect.Int64:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int64(i)
	case reflect.Uint:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int32(i)
	case reflect.Uint8:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint8(ui)
	case reflect.Uint16:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint16(ui)
	case reflect.Uint32:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint32(ui)
	case reflect.Uint64:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint64(ui)
	case reflect.Uintptr:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uintptr(ui)
	case reflect.Float32:
		str := string(data.([]uint8))
		f64, _ := strconv.ParseFloat(str, 64)
		data = float32(f64)
	case reflect.Float64:
		str := string(data.([]uint8))
		f64, _ := strconv.ParseFloat(str, 64)
		data = f64
	case reflect.Complex64:
		binBuf := bytes.NewBuffer(data.([]uint8))
		var x complex64
		_ = binary.Read(binBuf, binary.BigEndian, &x)
		data = x
	case reflect.Complex128:
		binBuf := bytes.NewBuffer(data.([]uint8))
		var x complex128
		_ = binary.Read(binBuf, binary.BigEndian, &x)
		data = x
	}

	return data
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
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = bytesToVal(data, tp)
			}
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
