package gobatis

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

func stringToVal(data interface{}, tp reflect.Type) interface{} {
	str := data.(string)
	switch tp.Kind() {
	case reflect.Bool:
		data = false
		if str == "1" {
			data = true
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
		f64, _ := strconv.ParseFloat(str, 64)
		data = float32(f64)
	case reflect.Float64:
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

func bytesToVal(data interface{}, tp reflect.Type) interface{} {
	str := string(data.([]uint8))
	switch tp.Kind() {
	case reflect.Bool:
		data = false
		if str == "1" {
			data = true
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
		f64, _ := strconv.ParseFloat(str, 64)
		data = float32(f64)
	case reflect.Float64:
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

func valToString(data interface{}) string {
	tp := reflect.TypeOf(data)
	s := ""
	switch tp.Kind() {
	case reflect.Bool:
		s = strconv.FormatBool(data.(bool))
	case reflect.Int:
		s = strconv.FormatInt(int64(data.(int)), 10)
	case reflect.Int8:
		s = strconv.FormatInt(int64(data.(int8)), 10)
	case reflect.Int16:
		s = strconv.FormatInt(int64(data.(int16)), 10)
	case reflect.Int32:
		s = strconv.FormatInt(int64(data.(int32)), 10)
	case reflect.Int64:
		s = strconv.FormatInt(int64(data.(int64)), 10)
	case reflect.Uint:
		s = strconv.FormatUint(uint64(data.(uint)), 10)
	case reflect.Uint8:
		s = strconv.FormatUint(uint64(data.(uint8)), 10)
	case reflect.Uint16:
		s = strconv.FormatUint(uint64(data.(uint16)), 10)
	case reflect.Uint32:
		s = strconv.FormatUint(uint64(data.(uint32)), 10)
	case reflect.Uint64:
		s = strconv.FormatUint(uint64(data.(uint64)), 10)
	case reflect.Uintptr:
		s = fmt.Sprint(data.(uintptr))
	case reflect.Float32:
		s = strconv.FormatFloat(float64(data.(float32)), 'f', -1, 64)
	case reflect.Float64:
		s = strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case reflect.Complex64:
		s = fmt.Sprint(data.(complex64))
	case reflect.Complex128:
		s = fmt.Sprint(data.(complex128))
	default:
		log.Println("[WARN]no process for type:" + tp.Name())
	}
	return s
}

func dataToFieldVal(data interface{}, tp reflect.Type, fieldName string) interface{} {
	defer func() {
		if err := recover(); nil != err {
			log.Println("[WARN] data to field val panic, fieldName:", fieldName, " err:", err)
		}
	}()

	typeName := tp.Name()
	switch {
	case typeName == "bool" ||
		typeName == "int" ||
		typeName == "int8" ||
		typeName == "int16" ||
		typeName == "int32" ||
		typeName == "int64" ||
		typeName == "uint" ||
		typeName == "uint8" ||
		typeName == "uint16" ||
		typeName == "uint32" ||
		typeName == "uint64" ||
		typeName == "uintptr" ||
		typeName == "float32" ||
		typeName == "float64" ||
		typeName == "complex64" ||
		typeName == "complex128":
		if nil != data {
			dataTp := reflect.TypeOf(data)
			if dataTp.Kind() == reflect.Slice ||
				dataTp.Kind() == reflect.Array {
				data = bytesToVal(data, tp)
			}

			dataTp = reflect.TypeOf(data)
			if dataTp.Kind() == reflect.String {
				data = stringToVal(data, tp)
			}

			return data
		}
	case typeName == "string":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				return string(data.([]byte))
			}

			data = valToString(data)
			return string(data.(string))
		}
	case typeName == "Time":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			tm, err := time.Parse("2006-01-02 15:04:05", data.(string))
			if err != nil {
				panic("time.Parse err:" + err.Error())
			}
			return tm
		}
	case typeName == "NullString":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}
			return NullString{String: data.(string), Valid: true}
		}
	case typeName == "NullInt64":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			i, err := strconv.ParseInt(data.(string), 10, 64)
			if err != nil {
				panic("ParseInt err:" + err.Error())
			}
			return NullInt64{Int64: i, Valid: true}
		}
	case typeName == "NullBool":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}
			if data.(string) == "true" {
				return NullBool{Bool: true, Valid: true}
			}
			return NullBool{Bool: false, Valid: true}
		}
	case typeName == "NullFloat64":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			f64, err := strconv.ParseFloat(data.(string), 64)
			if err != nil {
				panic("ParseFloat err:" + err.Error())
			}

			return NullFloat64{Float64: f64, Valid: true}
		}
	case typeName == "NullTime":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			tm, err := time.Parse("2006-01-02 15:04:05", data.(string))
			if err != nil {
				panic("time.Parse err:" + err.Error())
			}

			return NullTime{Time: tm, Valid: true}
		}
	}

	return nil
}
