package tagprocess

import (
	"github.com/wenj91/gobatis/tools/datautil"
	"github.com/wenj91/gobatis/tools/regutil"
	"github.com/wenj91/gobatis/xmlparser"
	"github.com/wenj91/gobatis/process/paramprocess"
	"errors"
	"reflect"
	"strings"
)

func ForeachProcess(foreachNode xmlparser.Node, param interface{}) (subSqlStr string, subSqlParams []interface{}, err error) {
	paramVal := reflect.ValueOf(param)
	kind := paramVal.Kind()
	finalParam := param

	//参数不能为指针
	if kind == reflect.Ptr {
		err = errors.New("param must not be ptr")
		return
	}

	//如果参数类型为struct,则取出对应名称集合
	if kind == reflect.Struct {
		list := foreachNode.Attr
		fieldName := list["collection"].Name.Local
		paramCollection := datautil.GetFieldValByName(param, fieldName)
		if nil == paramCollection {
			err = errors.New("no param match to foreach tag")
			return
		}
		finalParam = paramCollection
	}

	//如果参数类型为map,则取出对应字段名称集合
	if kind == reflect.Map {
		list := foreachNode.Attr
		paramCollection := param.(map[string]interface{})[list["collection"].Name.Local]
		if nil == paramCollection {
			err = errors.New("no param match to foreach tag")
			return
		}
		finalParam = paramCollection
	}

	//判断最后取出参数是否为数组
	finalParamVal := reflect.ValueOf(finalParam)
	finalKind := finalParamVal.Kind()
	if finalKind != reflect.Array && finalKind != reflect.Slice {
		err = errors.New("no param match to foreach tag")
		return
	}

	openStr := foreachNode.Attr["open"].Value
	closeStr := foreachNode.Attr["close"].Value
	separatorStr := foreachNode.Attr["separator"].Value

	subContent := ""
	eles := foreachNode.Elements
	for i := 0; i < len(eles); i++ {
		subContent += eles[i].Val.(string)
	}

	paramNames := make([]string, 0)

	//根据参数长度动态生成sql
	finalParamValLen := finalParamVal.Len()
	for i := 0; i < finalParamValLen; i++ {

		itemParamNames := make([]string, 0)
		if strings.Contains(subContent, "#{") {
			itemParamNames = regutil.SharpParamNamesFind(subContent)
			paramNames = append(paramNames, itemParamNames...)
		}

		//sql拼接处理
		if len(subSqlStr) > 0 {
			subSqlStr += separatorStr
		}
		subSqlStr += openStr
		subSqlStr += subContent
		subSqlStr += closeStr

		//参数提取
		//取出数组数据
		item := finalParamVal.Index(i)
		//判断元素类型, map, arr, struct
		paramVal := reflect.ValueOf(item)
		kind := paramVal.Kind()
		switch {
		//arr param process
		case kind == reflect.Array || kind == reflect.Slice:
			subSqlParams = append(subSqlParams, paramprocess.ArrayParamProcess(paramVal)...)
		case kind == reflect.Map: //map param process
			res, errTmp := paramprocess.MapParamProcess(item.Interface(), itemParamNames)
			if nil != errTmp {
				err = errTmp
				return
			}
			subSqlParams = append(subSqlParams, res...)
		case kind == reflect.Struct: //struct param process
			res, errTmp := paramprocess.StructParamProcess(item.Interface(), itemParamNames)
			if nil != errTmp {
				err = errTmp
				return
			}
			subSqlParams = append(subSqlParams, res...)
		}
	}

	//转化sql语句
	//paramNames
	paramsLen := len(paramNames)
	for i := 0; i < paramsLen/finalParamValLen; i++ {
		subSqlStr, err = regutil.SharpParamMatchReplace(subSqlStr, paramNames[i])
		if nil != err {
			return
		}
	}

	return
}
