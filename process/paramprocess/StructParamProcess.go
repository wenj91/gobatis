package paramprocess

import (
	"errors"
	"reflect"
	"github.com/wenj91/gobatis/tools/datautil"
)

// struct参数处理
func StructParamProcess(param interface{}, paramNames[]string)([]interface{}, error){

	sqlParams := make([]interface{}, 0)

	paramVal := reflect.ValueOf(param)
	if paramVal.Kind() == reflect.Ptr {
		return nil, errors.New("struct params must not be ptr")

	}
	for i := 0; i < len(paramNames); i++ {
		item := datautil.FieldToParams(param, paramNames[i])
		if nil == item {
			return nil, errors.New("no this params:" + paramNames[i])
		}
		sqlParams = append(sqlParams, item)
	}

	return sqlParams, nil
}
