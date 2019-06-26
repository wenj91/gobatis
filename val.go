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

func valUpcast(data interface{}, typeName string) interface{} {
	d := data
	switch typeName {
	case "bool":

	case "int":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = int(1)
			} else {
				d = int(0)
			}
		case int8:
			d = int(data.(int8))
		case int16:
			d = int(data.(int16))
		case int32:
			d = int(data.(int32))
		case int64:
			d = int(data.(int64))
		case uint:
			d = int(data.(uint))
		case uint8:
			d = int(data.(uint8))
		case uint16:
			d = int(data.(uint16))
		case uint32:
			d = int(data.(uint32))
		case uint64:
			d = int(data.(uint64))
		case uintptr:
			d = int(data.(uintptr))
		case float32:
			d = int(data.(float32))
		case float64:
			d = int(data.(float64))
		}
	case "int8":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = int8(1)
			} else {
				d = int8(0)
			}
		case int:
			d = int8(data.(int))
		case int16:
			d = int8(data.(int16))
		case int32:
			d = int8(data.(int32))
		case int64:
			d = int8(data.(int64))
		case uint:
			d = int8(data.(uint))
		case uint8:
			d = int8(data.(uint8))
		case uint16:
			d = int8(data.(uint16))
		case uint32:
			d = int8(data.(uint32))
		case uint64:
			d = int8(data.(uint64))
		case uintptr:
			d = int8(data.(uintptr))
		case float32:
			d = int8(data.(float32))
		case float64:
			d = int8(data.(float64))
		}
	case "int16":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = int16(1)
			} else {
				d = int16(0)
			}
		case int:
			d = int16(data.(int))
		case int8:
			d = int16(data.(int8))
		case int32:
			d = int16(data.(int32))
		case int64:
			d = int16(data.(int64))
		case uint:
			d = int16(data.(uint))
		case uint8:
			d = int16(data.(uint8))
		case uint16:
			d = int16(data.(uint16))
		case uint32:
			d = int16(data.(uint32))
		case uint64:
			d = int16(data.(uint64))
		case uintptr:
			d = int16(data.(uintptr))
		case float32:
			d = int16(data.(float32))
		case float64:
			d = int16(data.(float64))
		}
	case "int32":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = int32(1)
			} else {
				d = int32(0)
			}
		case int:
			d = int32(data.(int))
		case int8:
			d = int32(data.(int8))
		case int16:
			d = int32(data.(int16))
		case int64:
			d = int32(data.(int64))
		case uint:
			d = int32(data.(uint))
		case uint8:
			d = int32(data.(uint8))
		case uint16:
			d = int32(data.(uint16))
		case uint32:
			d = int32(data.(uint32))
		case uint64:
			d = int32(data.(uint64))
		case uintptr:
			d = int32(data.(uintptr))
		case float32:
			d = int32(data.(float32))
		case float64:
			d = int32(data.(float64))
		}
	case "int64":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = int64(1)
			} else {
				d = int64(0)
			}
		case int:
			d = int64(data.(int))
		case int8:
			d = int64(data.(int8))
		case int16:
			d = int64(data.(int16))
		case int32:
			d = int64(data.(int32))
		case uint:
			d = int64(data.(uint))
		case uint8:
			d = int64(data.(uint8))
		case uint16:
			d = int64(data.(uint16))
		case uint32:
			d = int64(data.(uint32))
		case uint64:
			d = int64(data.(uint64))
		case uintptr:
			d = int64(data.(uintptr))
		case float32:
			d = int64(data.(float32))
		case float64:
			d = int64(data.(float64))
		}
	case "uint":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = uint(1)
			} else {
				d = uint(0)
			}
		case int:
			d = uint(data.(int))
		case int8:
			d = uint(data.(int8))
		case int16:
			d = uint(data.(int16))
		case int32:
			d = uint(data.(int32))
		case int64:
			d = uint(data.(int64))
		case uint8:
			d = uint(data.(uint8))
		case uint16:
			d = uint(data.(uint16))
		case uint32:
			d = uint(data.(uint32))
		case uint64:
			d = uint(data.(uint64))
		case uintptr:
			d = uint(data.(uintptr))
		case float32:
			d = uint(data.(float32))
		case float64:
			d = uint(data.(float64))
		}
	case "uint8":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = uint8(1)
			} else {
				d = uint8(0)
			}
		case int:
			d = uint8(data.(int))
		case int8:
			d = uint8(data.(int8))
		case int16:
			d = uint8(data.(int16))
		case int32:
			d = uint8(data.(int32))
		case int64:
			d = uint8(data.(int64))
		case uint:
			d = uint8(data.(uint))
		case uint16:
			d = uint8(data.(uint16))
		case uint32:
			d = uint8(data.(uint32))
		case uint64:
			d = uint8(data.(uint64))
		case uintptr:
			d = uint8(data.(uintptr))
		case float32:
			d = uint8(data.(float32))
		case float64:
			d = uint8(data.(float64))
		}
	case "uint16":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = uint16(1)
			} else {
				d = uint16(0)
			}
		case int:
			d = uint16(data.(int))
		case int8:
			d = uint16(data.(int8))
		case int16:
			d = uint16(data.(int16))
		case int32:
			d = uint16(data.(int32))
		case int64:
			d = uint16(data.(int64))
		case uint:
			d = uint16(data.(uint))
		case uint8:
			d = uint16(data.(uint8))
		case uint32:
			d = uint16(data.(uint32))
		case uint64:
			d = uint16(data.(uint64))
		case uintptr:
			d = uint16(data.(uintptr))
		case float32:
			d = uint16(data.(float32))
		case float64:
			d = uint16(data.(float64))
		}
	case "uint32":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = uint32(1)
			} else {
				d = uint32(0)
			}
		case int:
			d = uint32(data.(int))
		case int8:
			d = uint32(data.(int8))
		case int16:
			d = uint32(data.(int16))
		case int32:
			d = uint32(data.(int32))
		case int64:
			d = uint32(data.(int64))
		case uint:
			d = uint32(data.(uint))
		case uint8:
			d = uint32(data.(uint8))
		case uint16:
			d = uint32(data.(uint16))
		case uint64:
			d = uint32(data.(uint64))
		case uintptr:
			d = uint32(data.(uintptr))
		case float32:
			d = uint32(data.(float32))
		case float64:
			d = uint32(data.(float64))
		}
	case "uint64":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = uint64(1)
			} else {
				d = uint64(0)
			}
		case int:
			d = uint64(data.(int))
		case int8:
			d = uint64(data.(int8))
		case int16:
			d = uint64(data.(int16))
		case int32:
			d = uint64(data.(int32))
		case int64:
			d = uint64(data.(int64))
		case uint:
			d = uint64(data.(uint))
		case uint8:
			d = uint64(data.(uint8))
		case uint16:
			d = uint64(data.(uint16))
		case uint32:
			d = uint64(data.(uint32))
		case uintptr:
			d = uint64(data.(uintptr))
		case float32:
			d = uint64(data.(float32))
		case float64:
			d = uint64(data.(float64))
		}
	case "uintptr":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = uintptr(1)
			} else {
				d = uintptr(0)
			}
		case int:
			d = uintptr(data.(int))
		case int8:
			d = uintptr(data.(int8))
		case int16:
			d = uintptr(data.(int16))
		case int32:
			d = uintptr(data.(int32))
		case int64:
			d = uintptr(data.(int64))
		case uint:
			d = uintptr(data.(uint))
		case uint8:
			d = uintptr(data.(uint8))
		case uint16:
			d = uintptr(data.(uint16))
		case uint32:
			d = uintptr(data.(uint32))
		case uint64:
			d = uintptr(data.(uint64))
		case float32:
			d = uintptr(data.(float32))
		case float64:
			d = uintptr(data.(float64))
		}
	case "float32":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = float32(1)
			} else {
				d = float32(0)
			}
		case int:
			d = float32(data.(int))
		case int8:
			d = float32(data.(int8))
		case int16:
			d = float32(data.(int16))
		case int32:
			d = float32(data.(int32))
		case int64:
			d = float32(data.(int64))
		case uint:
			d = float32(data.(uint))
		case uint8:
			d = float32(data.(uint8))
		case uint16:
			d = float32(data.(uint16))
		case uint32:
			d = float32(data.(uint32))
		case uint64:
			d = float32(data.(uint64))
		case uintptr:
			d = float32(data.(uintptr))
		case float64:
			d = float32(data.(float64))
		}
	case "float64":
		switch data.(type) {
		case bool:
			if data.(bool) {
				d = float64(1)
			} else {
				d = float64(0)
			}
		case int:
			d = float64(data.(int))
		case int8:
			d = float64(data.(int8))
		case int16:
			d = float64(data.(int16))
		case int32:
			d = float64(data.(int32))
		case int64:
			d = float64(data.(int64))
		case uint:
			d = float64(data.(uint))
		case uint8:
			d = float64(data.(uint8))
		case uint16:
			d = float64(data.(uint16))
		case uint32:
			d = float64(data.(uint32))
		case uint64:
			d = float64(data.(uint64))
		case uintptr:
			d = float64(data.(uintptr))
		case float32:
			d = float64(data.(float32))
		}
	case "complex64":
	case "complex128":

	}

	return d
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

			data = valUpcast(data, typeName)

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
