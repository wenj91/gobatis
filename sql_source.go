package gobatis

import (
	"log"
	"strings"
)

// gobatis的核心, 从配置到sql, 参数映射......
type boundSql struct {
	sqlStr        string
	paramMappings []string
	params        map[string]interface{}
	extParams     map[string]interface{}
}

type iSqlSource interface {
	getBoundSql(params map[string]interface{}) *boundSql
}

type dynamicSqlSource struct {
	sqlNode iSqlNode
}

func (this *dynamicSqlSource) getBoundSql(params map[string]interface{}) *boundSql {
	ctx := &dynamicContext{params: params}
	this.sqlNode.build(ctx)

	sss := staticSqlSource{
		sqlStr: ctx.sqlStr,
	}

	bs := sss.getBoundSql(params)
	bs.extParams = ctx.params

	return bs
}

type staticSqlSource struct {
	sqlStr        string
	paramMappings []string
}

func (this *staticSqlSource) getBoundSql(params map[string]interface{}) *boundSql {
	this.dollarTokenHandler(params)
	this.tokenHandler(params)
	return &boundSql{
		sqlStr:        this.sqlStr,
		paramMappings: this.paramMappings,
		params:        params,
	}
}

// ${xx}处理
func (this *staticSqlSource) dollarTokenHandler(params map[string]interface{}) {
	sqlStr := this.sqlStr
	if strings.Index(sqlStr, "$") == -1 {
		return
	}

	finalSqlStr := ""
	itemStr := ""
	start := 0
	for i := 0; i < len(sqlStr); i++ {
		if start > 0 {
			itemStr += string(sqlStr[i])
		}

		if i != 0 && i < len(sqlStr) {
			if string([]byte{sqlStr[i-1], sqlStr[i]}) == "${" {
				start = i
			}
		}

		if start != 0 && i < len(sqlStr)-1 && sqlStr[i+1] == '}' {
			finalSqlStr += sqlStr[:start-1]
			sqlStr = sqlStr[i+2:]

			itemStr = strings.Trim(itemStr, " ")
			//this.paramMappings = append(this.paramMappings, itemStr)

			item, ok := params[itemStr]
			if !ok {
				log.Println("ERROR: param ${" + itemStr + "}, not found")
				panic("params:" + itemStr + " not found")
			}

			finalSqlStr += item.(string)

			i = 0
			start = 0
			itemStr = ""
		}
	}

	if start != 0 {
		log.Println("WARN: token not close")
	}

	finalSqlStr += sqlStr
	finalSqlStr = strings.Trim(finalSqlStr, " ")
	this.sqlStr = finalSqlStr
}

// 静态token处理, 将#{xx}预处理为数据库预编译语句
func (this *staticSqlSource) tokenHandler(params map[string]interface{}) {
	sqlStr := this.sqlStr

	finalSqlStr := ""
	itemStr := ""
	start := 0
	for i := 0; i < len(sqlStr); i++ {
		if start > 0 {
			itemStr += string(sqlStr[i])
		}

		if i != 0 && i < len(sqlStr) {
			if string([]byte{sqlStr[i-1], sqlStr[i]}) == "#{" {
				start = i
			}
		}

		if start != 0 && i < len(sqlStr)-1 && sqlStr[i+1] == '}' {
			finalSqlStr += sqlStr[:start-1]
			sqlStr = sqlStr[i+2:]

			itemStr = strings.Trim(itemStr, " ")
			this.paramMappings = append(this.paramMappings, itemStr)

			finalSqlStr += "?"

			i = 0
			start = 0
			itemStr = ""
		}
	}

	if start != 0 {
		log.Println("WARN: token not close")
	}

	finalSqlStr += sqlStr
	finalSqlStr = strings.Trim(finalSqlStr, " ")
	this.sqlStr = finalSqlStr
}
