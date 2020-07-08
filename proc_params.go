package gobatis

import (
	"reflect"
	"strconv"
	"time"
)

// parameters process util
// @params
//    param interface{} : sql query params
// @return
//    map[string]interface{} : return the convert map
func paramProcess(param interface{}) map[string]interface{} {
	v := reflect.ValueOf(param)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	res := make(map[string]interface{})
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		LOG.Warn("Foreach tag collection element must not be slice or array")
		res = listToMap(param)
	case reflect.Struct:
		res = structToMap(param)
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
func listToMap(arr interface{}) map[string]interface{} {
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
func structToMap(s interface{}) map[string]interface{} {
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
			ns := s.(NullString)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullInt64":
		res["0"] = nil
		if nil != s {
			ns := s.(NullInt64)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullBool":
		res["0"] = nil
		if nil != s {
			ns := s.(NullBool)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullFloat64":
		res["0"] = nil
		if nil != s {
			ns := s.(NullFloat64)
			if ns.Valid {
				str, _ := ns.Value()
				res["0"] = str
			}
		}
	case "NullTime":
		res["0"] = nil
		if nil != s {
			ns := s.(NullTime)
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

				data := fieldToVal(fieldVal.Interface())
				res[field.Name] = data
				// 同时可以使用tag做参数名 https://github.com/wenj91/gobatis/issues/43
				tag := field.Tag.Get("field")
				if tag != "" && tag != "-" {
					res[tag] = data
				}
			}
		}
	}

	return res
}

func fieldToVal(field interface{}) interface{} {
	objVal := reflect.ValueOf(field)
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	tp := objVal.Type()
	switch tp.Name() {
	case "Time":
		if nil != field {
			return field.(time.Time).Format("2006-01-02 15:04:05")
		}
	case "NullString":
		if nil != field {
			ns := field.(NullString)
			if ns.Valid {
				str, _ := ns.Value()
				return str
			}
		}
	case "NullInt64":
		if nil != field {
			ni64 := field.(NullInt64)
			if ni64.Valid {
				i, _ := ni64.Value()
				return i
			}
		}
	case "NullBool":
		if nil != field {
			nb := field.(NullBool)
			if nb.Valid {
				b, _ := nb.Value()
				return b
			}
		}
	case "NullFloat64":
		if nil != field {
			nf := field.(NullFloat64)
			if nf.Valid {
				f, _ := nf.Value()
				return f
			}
		}
	case "NullTime":
		if nil != field {
			nt := field.(NullTime)
			if nt.Valid {
				t, _ := nt.Value()
				return t.(time.Time).Format("2006-01-02 15:04:05")
			}
		}
	default:
		return field
	}

	return nil
}
