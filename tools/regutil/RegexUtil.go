package regutil

import (
	"regexp"
	"strings"
	"errors"
)

const (
	SQL_PARAM_REGEX_SHARP       = `#{\s*id\s*}`               //匹配#{xxx}类型正则表达式,用于替换规则
	SQL_PARAM_REGEX_DOLLAR      = `\${\s*id\s*}`              //匹配${xxx}类型正则表达式
	SQL_PARAM_REGEX_SHARP_MATCH = `#{\s*[A-Za-z0-9_$,=]+\s*}` //匹配#{xxx}类型正则表达式,用于查找规则
)

var NO_PARAM_MATCH_ERR = errors.New("no param match!")

// 删除str中的空格
func SpaceMatchReplace(str string) string {
	reg := regexp.MustCompile("\\s*")
	return string(reg.ReplaceAll([]byte(str), []byte("")))
}

// 查找sqlTemplate中参数集合
func SharpParamNamesFind(sqlTemplate string) []string {
	reg := regexp.MustCompile(SQL_PARAM_REGEX_SHARP_MATCH)
	paramNames := reg.FindAllString(sqlTemplate, -1)
	for i := 0; i < len(paramNames); i++ {
		paramNames[i] = strings.Replace(
			strings.Replace(SpaceMatchReplace(paramNames[i]), "#{", "", -1),
			"}", "", -1)
	}

	return paramNames
}

// 将sqlTemplate中的参数集合全部替换为?
func SharpParamMatchReplace(sqlTemplate string, paramName string) (string, error) {
	return SharpParamMatchReplaceByVal(sqlTemplate, paramName, " ? ")
}

// 将sqlTemplate中的参数集合全部替换为指定val
func SharpParamMatchReplaceByVal(sqlTemplate string, paramName string, val string) (string, error) {
	regexStr := strings.Replace(SQL_PARAM_REGEX_SHARP, "id", paramName, -1)

	reg := regexp.MustCompile(regexStr)
	if !reg.Match([]byte(sqlTemplate)) {
		return "", NO_PARAM_MATCH_ERR
	}

	res := reg.ReplaceAll([]byte(sqlTemplate), []byte(val))

	return string(res), nil
}

func DollarParamMatchReplace(sqlTemplate string, paramName string, val string) (string, error) {
	regexStr := strings.Replace(SQL_PARAM_REGEX_DOLLAR, "id", paramName, -1)

	reg := regexp.MustCompile(regexStr)
	if !reg.Match([]byte(sqlTemplate)) {
		return "", NO_PARAM_MATCH_ERR
	}

	res := reg.ReplaceAll([]byte(sqlTemplate), []byte(val))

	return string(res), nil
}
