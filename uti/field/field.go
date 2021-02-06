package field

import (
	"github.com/wenj91/gobatis/uti/param"
	"reflect"
)

// Fields return model sql mapper fields
func Fields(m interface{}) []string {
	li := make([]string, 0)
	if nil != m {
		objType := reflect.TypeOf(m)
		if objType.Kind() == reflect.Ptr {
			objType = objType.Elem()
		}

		for i := 0; i < objType.NumField(); i++ {
			field := objType.Field(i)
			tag := field.Tag.Get("field")
			if tag != "" && tag != "-" {
				li = append(li, tag)
			}
		}
	}

	return li
}

// Map return model sql mapper fields
func Map(m interface{}) (cols []string, vals []interface{}) {
	if nil != m {
		objVal := reflect.ValueOf(m)
		objType := reflect.TypeOf(m)
		if objType.Kind() == reflect.Ptr {
			objType = objType.Elem()
			objVal = objVal.Elem()
		}

		for i := 0; i < objType.NumField(); i++ {
			field := objVal.Field(i)
			k := field.Kind()
			if k == reflect.Chan ||
				k == reflect.Func ||
				k == reflect.Map ||
				k == reflect.Ptr ||
				k == reflect.UnsafePointer ||
				k == reflect.Interface ||
				k == reflect.Slice {
				if !field.IsNil() {
					ft := objType.Field(i)
					tag := ft.Tag.Get("field")
					if tag != "" && tag != "-" {
						f, _ := param.FieldToVal(field.Interface())
						cols = append(cols, tag)
						vals = append(vals, f)
					}
				}
			} else {
				ft := objType.Field(i)
				tag := ft.Tag.Get("field")
				if tag != "" && tag != "-" {
					cols = append(cols, tag)
					vals = append(vals, field.Interface())
				}
			}
		}
	}

	return
}
