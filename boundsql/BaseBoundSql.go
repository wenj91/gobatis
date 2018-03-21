package boundsql

import (
	"github.com/wenj91/gobatis/process/paramprocess"
	"github.com/wenj91/gobatis/tools/regutil"
	"errors"
	"reflect"
	"strings"
)

//sql基本处理
func BaseGetBoundSql(sqlStr string, param interface{}) (boundSql BoundSql, err error) {

	paramNames := []string{}
	if strings.Contains(sqlStr, "#{") {
		paramNames = regutil.SharpParamNamesFind(sqlStr)
	} else {
		//如果没有标识,说明不用处理参数
		boundSql.Sql = sqlStr
		return
	}

	//转化sql语句
	for i := 0; i < len(paramNames); i++ {
		sqlStr, err = regutil.SharpParamMatchReplace(sqlStr, paramNames[i])
		if nil != err {
			return
		}
	}

	paramVal := reflect.ValueOf(param)
	kind := paramVal.Kind()
	sqlParams := make([]interface{}, 0)
	switch {
	//arr param process
	case kind == reflect.Array || kind == reflect.Slice:
		sqlParams = append(sqlParams, paramprocess.ArrayParamProcess(paramVal)...)
	case kind == reflect.Map: //map param process
		res, errTmp := paramprocess.MapParamProcess(param, paramNames)
		if nil != errTmp {
			err = errTmp
			return
		}
		sqlParams = append(sqlParams, res...)
	case kind == reflect.Struct: //struct param process
		res, errTmp := paramprocess.StructParamProcess(param, paramNames)
		if nil != errTmp {
			err = errTmp
			return
		}
		sqlParams = append(sqlParams, res...)
	case kind == reflect.Bool ||
		kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Uintptr ||
		kind == reflect.Float32 ||
		kind == reflect.Float64 ||
		kind == reflect.Complex64 ||
		kind == reflect.Complex128 ||
		kind == reflect.String: //base type param process
		sqlParams = append(sqlParams, param)
	}

	if len(paramNames) != len(sqlParams) {
		err = errors.New("param size error")
		return
	}

	boundSql.ParameterMappings = sqlParams
	boundSql.Sql = sqlStr

	return
}
