package resultprocess

import (
	"log"
	"database/sql"
	"reflect"
	"errors"
)

//处理结果生成[]map[string][]byte
func MapsProcess(rows *sql.Rows, result interface{}, sqlParams []interface{})(int, error) {
	resPtr := reflect.ValueOf(result)
	if resPtr.Kind() != reflect.Ptr {
		return 0, errors.New("Maps query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return 0, errors.New("Maps query result must be slice prt")
	}

	count := 0
	for rows.Next() {
		count++
		resMap := make(map[string]interface{})
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
			resMap[cols[i]] = vals[i]
		}
		value.Set(reflect.Append(value, reflect.ValueOf(resMap)))
	}

	return count, nil
}
