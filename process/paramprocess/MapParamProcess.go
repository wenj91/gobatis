package paramprocess

import "errors"

// map参数处理
func MapParamProcess(param interface{}, paramNames[]string)([]interface{}, error){

	sqlParams := make([]interface{}, 0)

	paramMap := param.(map[string]interface{})
	for i := 0; i < len(paramNames); i++ {
		item := paramMap[paramNames[i]]
		if nil == item {
			return nil, errors.New("params must not be nil")
		}
		sqlParams = append(sqlParams, item)
	}

	return sqlParams, nil
}
