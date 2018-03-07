package resultprocess

import (
	"log"
	"database/sql"
	"reflect"
	"errors"
)

//处理结果生成map[string][]byte
func MapProcess(rows *sql.Rows, result interface{}, sqlParams []interface{})(int, error) {

	resBean := reflect.ValueOf(result)
	if resBean.Kind() == reflect.Ptr{
		return 0, errors.New("Map query result can not be ptr")
	}

	//id,email,head_image_url
	count := 0
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
			result.(map[string]interface{})[cols[i]] = vals[i]
		}
	}

	return 1, nil
}
