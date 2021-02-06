package param

import (
	"github.com/wenj91/gobatis/logger"
	"github.com/wenj91/gobatis/na"
	"github.com/wenj91/gobatis/uti"
	"reflect"
	"strconv"
	"time"
)

// parameters process util
// @params
//    param interface{} : sql query params
// @return
//    map[string]interface{} : return the convert map
func Process(param interface{}) map[string]interface{} {
	v := reflect.ValueOf(param)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	res := make(map[string]interface{})
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		logger.LOG.Warn("Foreach tag collection element must not be slice or array")
		res = ListToMap(param)
	case reflect.Struct:
		res = StructToMap(param)
	case reflect.Map:
		res = param.(map[string]interface{})
	default:
		res["0"] = param
	}

	return res
}

// convert list to map
// @params
//    arr interface{} : list param
// @return
//    map[string]interface{} : return the convert map
func ListToMap(arr interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	objVal := reflect.ValueOf(arr)
	if objVal.Kind() != reflect.Array && objVal.Kind() != reflect.Slice {
		return res
	}

	res["list"] = arr

	for i := 0; i < objVal.Len(); i++ {
		res[strconv.Itoa(i)] = objVal.Index(i).Interface()
	}

	return res
}

// convert struct to map
// @params
//    s interface{} : struct param
// @return
//    map[string]interface{} : return the convert map
func StructToMap(s interface{}) map[string]interface{} {
	objVal := reflect.ValueOf(s)
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	res := make(map[string]interface{})

	tp := objVal.Type()
	switch tp.Name() {
	case "Time":
		res["0"] = nil
		if nil != s {
			res["0"] = s.(time.Time).Format("2006-01-02 15:04:05")
		}
	case "NullString":
		res["0"] = nil
		if nil != s {
			ns := s.(na.NullString)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullInt64":
		res["0"] = nil
		if nil != s {
			ns := s.(na.NullInt64)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullBool":
		res["0"] = nil
		if nil != s {
			ns := s.(na.NullBool)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullFloat64":
		res["0"] = nil
		if nil != s {
			ns := s.(na.NullFloat64)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullTime":
		res["0"] = nil
		if nil != s {
			ns := s.(na.NullTime)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	default:
		objType := objVal.Type()
		for i := 0; i < objVal.NumField(); i++ {
			fieldVal := objVal.Field(i)
			if fieldVal.CanInterface() {
				field := objType.Field(i)

				data, ok := fieldToVal(fieldVal.Interface())
				if ok {
					res[field.Name] = data
					// 同时可以使用tag做参数名 https://github.com/wenj91/gobatis/issues/43
					tag := field.Tag.Get("field")
					if tag != "" && tag != "-" {
						res[tag] = data
					}
				}
			}
		}
	}

	return res
}

func fieldToVal(field interface{}) (interface{}, bool) {
	isNil, _ := uti.IsNil(field)
	if !isNil {
		return nil, false
	}

	objVal := reflect.ValueOf(field)
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	tp := objVal.Type()
	switch tp.Name() {
	case "Time":
		return field.(time.Time).Format("2006-01-02 15:04:05"), true
	case "NullString":
		ns := field.(na.NullString)
		if ns.Valid {
			str, _ := ns.Value()
			return str, true
		}
	case "NullInt64":
		ni64 := field.(na.NullInt64)
		if ni64.Valid {
			i, _ := ni64.Value()
			return i, true
		}
	case "NullBool":
		nb := field.(na.NullBool)
		if nb.Valid {
			b, _ := nb.Value()
			return b, true
		}
	case "NullFloat64":
		nf := field.(na.NullFloat64)
		if nf.Valid {
			f, _ := nf.Value()
			return f, true
		}
	case "NullTime":
		nt := field.(na.NullTime)
		if nt.Valid {
			t, _ := nt.Value()
			return t.(time.Time).Format("2006-01-02 15:04:05"), true
		}
	default:
		return field, true
	}

	return nil, false
}
