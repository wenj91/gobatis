package gobatis

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

type dynamicContext struct {
	sqlStr string
	params map[string]interface{}
}

func (this *dynamicContext) appendSql(sqlStr string) {
	this.sqlStr += sqlStr + " "
}

// [ref](http://www.mybatis.org/mybatis-3/dynamic-sql.html)
type iSqlNode interface {
	build(ctx *dynamicContext) bool
}

// mixed node
type mixedSqlNode struct {
	sqlNodes []iSqlNode
}

func (this *mixedSqlNode) build(ctx *dynamicContext) bool {
	for i := 0; i < len(this.sqlNodes); i++ {
		sqlNode := this.sqlNodes[i]
		sqlNode.build(ctx)
	}

	return true
}

// if node
type ifSqlNode struct {
	test    string
	sqlNode iSqlNode
}

func (this *ifSqlNode) build(ctx *dynamicContext) bool {
	if ok := eval(this.test, ctx.params); ok {
		this.sqlNode.build(ctx)
		return true
	}

	return false
}

// text node
type textSqlNode struct {
	content string
}

func (this *textSqlNode) build(ctx *dynamicContext) bool {
	ctx.appendSql(this.content)
	return true
}

// for node
const listItemPrefix = "_ls_item_p_"

type foreachSqlNode struct {
	sqlNode    iSqlNode
	collection string
	open       string
	close      string
	separator  string
	item       string
	index      string
}

func (this *foreachSqlNode) build(ctx *dynamicContext) bool {
	collection, ok := ctx.params[this.collection]
	if !ok {
		log.Println("No collection for foreach tag:", this.collection)
		return false
	}

	ctx.appendSql(this.open)

	val := reflect.ValueOf(collection)

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		log.Println("Foreach tag collection must be slice or array")
		return false
	}

	for i := 0; i < val.Len(); i++ {
		v := val.Index(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// convert struct map val to params
		keys := make([]string, 0)
		params := make(map[string]interface{})
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			log.Println("Foreach tag collection element must not be slice or array")
			return false
		case reflect.Struct:
			m := this.structToMap(v.Interface())
			for k, v := range m {
				key := this.item + "." + k
				keys = append(keys, key)
				params[key] = v
			}
		case reflect.Map:
			m := v.Interface().(map[string]interface{})
			for k, v := range m {
				key := this.item + "." + k
				keys = append(keys, key)
				params[key] = v
			}
		default:
			keys = append(keys, this.item)
			params[this.item] = v.Interface()
		}

		params[this.item] = v.Interface()

		tempCtx := &dynamicContext{
			params: params,
		}

		this.sqlNode.build(tempCtx)
		this.tokenHandler(tempCtx, i)

		if i != 0 {
			ctx.appendSql(this.separator)
		}

		ctx.appendSql(tempCtx.sqlStr)

		// del temp param
		for _, k := range keys {
			delete(tempCtx.params, k)
		}

		// sync tempCtx params to ctx
		for k, v := range tempCtx.params {
			ctx.params[k] = v
		}
	}
	ctx.appendSql(this.close)

	return true
}

func (this *foreachSqlNode) tokenHandler(ctx *dynamicContext, index int) {
	sqlStr := ctx.sqlStr

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
			finalSqlStr += sqlStr[:start+1]
			sqlStr = sqlStr[i+2:]

			var re = regexp.MustCompile("^\\s*" + this.item + "\\s*")
			itemPrefix := listItemPrefix + this.item + fmt.Sprintf("%d", index)
			s := re.ReplaceAllString(itemStr, itemPrefix)
			s = strings.Trim(s, " ")
			if strings.Contains(s, itemPrefix) {
				itemKey := strings.Trim(itemStr, " ")
				if v, ok := ctx.params[itemKey]; ok {
					ctx.params[s] = v
				}
			}

			finalSqlStr += s + "}"

			i = 0
			start = 0
			itemStr = ""
		}
	}

	if start != 0 {
		log.Println("WARN: token not close, SqlStr:" + ctx.sqlStr + " At:" + fmt.Sprintf("%d", start))
	}

	finalSqlStr += sqlStr
	ctx.sqlStr = finalSqlStr
}

func (this *foreachSqlNode) structToMap(s interface{}) map[string]interface{} {
	return structToMap(s)
}

// set node
type setSqlNode struct {
	sqlNodes []iSqlNode
}

func (this *setSqlNode) build(ctx *dynamicContext) bool {

	sqlStr := ""
	for _, sqlNode := range this.sqlNodes {
		tempCtx := &dynamicContext{
			params: ctx.params,
		}
		sqlNode.build(tempCtx)
		if sqlStr != "" && tempCtx.sqlStr != "" {
			sqlStr += " , "
		}

		sqlStr += tempCtx.sqlStr

		for k, v := range tempCtx.params {
			ctx.params[k] = v
		}
	}

	if sqlStr != "" {
		ctx.appendSql(" set ")
		sqlStr = strings.TrimSpace(sqlStr)
		sqlStr = strings.TrimSuffix(sqlStr, ",")
		ctx.appendSql(sqlStr)
	}

	return true
}

// trim node
type trimSqlNode struct {
	prefix          string // prefix：前缀
	prefixOverrides string // prefixOverride：去掉第一个出现prefixOverrides字符串
	suffixOverrides string // suffixOverride：去掉最后一个字符串
	suffix          string // suffix：后缀
	sqlNodes        []iSqlNode
}

func (this *trimSqlNode) build(ctx *dynamicContext) bool {
	tempCtx := &dynamicContext{
		params: ctx.params,
	}

	for _, sqlNode := range this.sqlNodes {
		if tempCtx.sqlStr != "" {
			tempCtx.sqlStr += " "
		}
		sqlNode.build(tempCtx)
	}

	if tempCtx.sqlStr != "" {
		sqlStr := strings.TrimSpace(tempCtx.sqlStr)

		preOv := strings.TrimSpace(this.prefixOverrides)
		if preOv != "" {
			sqlStr = strings.TrimPrefix(sqlStr, preOv+" ")
		}

		suffOv := strings.TrimSpace(this.suffixOverrides)
		if suffOv != "" {
			sqlStr = strings.TrimSuffix(sqlStr, suffOv+" ")
		}

		pre := strings.TrimSpace(this.prefix)
		if pre != "" {
			sqlStr = pre + " " + sqlStr
		}

		suff := strings.TrimSpace(this.suffix)
		if suff != "" {
			sqlStr += " " + suff
		}

		ctx.appendSql(sqlStr)
	}

	for k, v := range tempCtx.params {
		ctx.params[k] = v
	}

	return true
}

// where node
type whereSqlNode struct {
	sqlNodes []iSqlNode
}

func (this *whereSqlNode) build(ctx *dynamicContext) bool {
	tempCtx := &dynamicContext{
		params: ctx.params,
	}

	for _, sqlNode := range this.sqlNodes {
		if tempCtx.sqlStr != "" {
			tempCtx.sqlStr += " "
		}
		sqlNode.build(tempCtx)
	}

	if tempCtx.sqlStr != "" {
		sqlStr := strings.TrimSpace(tempCtx.sqlStr)
		sqlStr = strings.TrimPrefix(sqlStr, "and ")
		sqlStr = strings.TrimPrefix(sqlStr, "AND ")
		sqlStr = strings.TrimPrefix(sqlStr, "or ")
		sqlStr = strings.TrimPrefix(sqlStr, "OR ")

		ctx.appendSql(" where ")
		ctx.appendSql(sqlStr)
	}

	for k, v := range tempCtx.params {
		ctx.params[k] = v
	}

	return true
}

// choose node
