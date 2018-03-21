package paramprocess

import "reflect"

// 数组参数处理
func ArrayParamProcess(paramVal reflect.Value) []interface{}{
	paramLen := paramVal.Len()

	sqlParams := make([]interface{}, 0)
	for i := 0; i < paramLen; i++ {
		itemVal := paramVal.Index(i)
		//数组参数校验,如果数据元素为struct, 如果不是nullbool, nullfloat64, nullint64, nullstring, nulltime, 则返回异常
		sqlParams = append(sqlParams, itemVal.Interface())
	}

	return sqlParams
}
