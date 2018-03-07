package resultprocess

import (
	"log"
	"database/sql"
	"reflect"
	"errors"
	"gobatis/tools/datautil"
)

func StructsProcess(rows *sql.Rows, result interface{}, sqlParams []interface{})(int, error) {

	sliceVal := reflect.ValueOf(result)
	if sliceVal.Kind() != reflect.Ptr {
		return 0, errors.New("Structs query result must be ptr")
	}

	slicePtr := reflect.Indirect(sliceVal)
	if slicePtr.Kind() != reflect.Slice {
		return 0, errors.New("Structs query result must be slice")
	}

	//get elem type
	elem := slicePtr.Type().Elem()
	resultType := elem
	isPtr := (elem.Kind() == reflect.Ptr)
	if isPtr {
		resultType = elem.Elem()
	}

	if resultType.Kind() != reflect.Struct {
		return 0, errors.New("Structs query results item  must be struct")
	}

	count := 0
	for rows.Next() {
		count++
		cols, err := rows.Columns()
		if nil != err {
			log.Fatal(err)
			return 0, err
		}

		vals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}

		rows.Scan(scanArgs...)

		obj := reflect.New(resultType).Elem()

		objPtr := reflect.Indirect(obj)
		for i := 0; i < len(cols); i++ {
			field := objPtr.FieldByName(cols[i])
			//设置相关字段的值,并判断是否可设值
			if field.CanSet() && vals[i] != nil {

				//获取字段类型并设值
				data := datautil.DataToFieldVal(vals[i], field.Type())
				if nil != data {
					field.Set(reflect.ValueOf(data))
				}
			}
		}

		if isPtr {
			slicePtr.Set(reflect.Append(slicePtr, obj.Addr()))
		} else {
			slicePtr.Set(reflect.Append(slicePtr, obj))
		}
	}

	return count, nil
}
