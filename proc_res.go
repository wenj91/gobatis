package gobatis

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

type resultTypeProc = func(rows *sql.Rows, res interface{}) error

var resSetProcMap = map[ResultType]resultTypeProc{
	resultTypeMap:    resMapProc,
	resultTypeMaps:   resMapsProc,
	resultTypeSlice:  resSliceProc,
	resultTypeSlices: resSlicesProc,
	resultTypeValue:  resValueProc,
}

func resValueProc(rows *sql.Rows, res interface{}) error {

	return nil
}

func resSlicesProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("Slices query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return errors.New("Slices query result must be slice ptr")
	}

	arr, err := rowsToSlices(rows)
	if nil != err {
		return err
	}

	for i := 0; i < len(arr); i++ {
		value.Set(reflect.Append(value, reflect.ValueOf(arr[i])))
	}

	return nil
}

func resSliceProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("Slice query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return errors.New("Slice query result must be slice ptr")
	}

	arr, err := rowsToSlices(rows)
	if nil != err {
		return err
	}

	if len(arr) > 1 {
		return errors.New("Slice query result more than one row")
	}

	if len(arr) > 0 {
		tempResSlice := arr[0].([]interface{})
		value.Set(reflect.AppendSlice(value, reflect.ValueOf(tempResSlice)))
	}

	return nil
}

func resMapProc(rows *sql.Rows, res interface{}) error {
	resBean := reflect.ValueOf(res)
	if resBean.Kind() == reflect.Ptr {
		return errors.New("Map query result can not be ptr")
	}

	if resBean.Kind() != reflect.Map {
		return errors.New("Map query result must be map")
	}

	arr, err := rowsToMaps(rows)
	if nil != err {
		return err
	}

	if len(arr) > 1 {
		return errors.New("Map query result more than one row")
	}

	if len(arr) > 0 {
		resMap := res.(map[string]interface{})
		tempResMap := arr[0].(map[string]interface{})
		for k, v := range tempResMap {
			resMap[k] = v
		}
	}

	return nil
}

func resMapsProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("Maps query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return errors.New("Maps query result must be slice ptr")
	}
	arr, err := rowsToMaps(rows)
	if nil != err {
		return err
	}

	for i := 0; i < len(arr); i++ {
		value.Set(reflect.Append(value, reflect.ValueOf(arr[i])))
	}

	return nil
}

func rowsToMaps(rows *sql.Rows) ([]interface{}, error) {
	res := make([]interface{}, 0)
	for rows.Next() {
		resMap := make(map[string]interface{})
		cols, err := rows.Columns()
		if nil != err {
			log.Println(err)
			return res, err
		}

		vals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}

		rows.Scan(scanArgs...)
		for i := 0; i < len(cols); i++ {
			val := vals[i]
			if nil != val {
				v := reflect.ValueOf(val)
				if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
					val = string(val.([]uint8))
				}
			}
			resMap[cols[i]] = val
		}

		res = append(res, resMap)
	}

	return res, nil
}

func rowsToSlices(rows *sql.Rows) ([]interface{}, error) {
	res := make([]interface{}, 0)
	for rows.Next() {
		resSlice := make([]interface{}, 0)
		cols, err := rows.Columns()
		if nil != err {
			log.Println(err)
			return nil, err
		}

		vals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}

		rows.Scan(scanArgs...)
		for i := 0; i < len(cols); i++ {
			val := vals[i]
			if nil != val {
				v := reflect.ValueOf(val)
				if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
					val = string(val.([]uint8))
				}
			}
			resSlice = append(resSlice, val)
		}

		res = append(res, resSlice)
	}

	return res, nil
}
