package resultprocess

import (
	"log"
	"database/sql"
	"reflect"
	"errors"
	"gobatis/tools/datautil"
)

func StructProcess(rows *sql.Rows, result interface{}, sqlParams []interface{})(int, error) {

	resBean := reflect.ValueOf(result)
	if resBean.Kind() != reflect.Ptr{
		return 0, errors.New("Struct query result must be ptr")
	}

	//id,email,head_image_url
	count := 0
	for rows.Next() {
		count++
		if count>1{
			return 0, errors.New("Struct query return more than 1 result")
		}
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
		objPtr := reflect.Indirect(resBean)
		for i := 0; i < len(cols); i++ {
			field := objPtr.FieldByName(cols[i])
			//设置相关字段的值
			if field.CanSet() && vals[i] != nil {

				//获取字段类型并设值
				data := datautil.DataToFieldVal(vals[i], field.Type())
				if nil != data {
					field.Set(reflect.ValueOf(data))
				}
			}
		}
	}

	return count, nil
}
