package resultprocess

import (
	"log"
	"database/sql"
	"reflect"
	"errors"
)

//处理结果生成[][]byte
func SliceProcess(rows *sql.Rows, result interface{}, sqlParams []interface{})(int, error) {

	resPtr := reflect.ValueOf(result)
	if resPtr.Kind() != reflect.Ptr {
		return 0, errors.New("Maps query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return 0, errors.New("Maps query result must be slice prt")
	}

	//id,email,head_image_url
	count := 0
	res := []interface{}{}
	for rows.Next() {
		count++
		if count>1{
			return 0, errors.New("Map query return more than 1 result")
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
		for i := 0; i < len(cols); i++ {
			res = append(res, vals[i])
		}
	}

	value.Set(reflect.ValueOf(res))

	return 1, nil
}
