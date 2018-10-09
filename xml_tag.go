package gobatis

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type dynamicContext struct {
	sqlStr string
	params map[string]interface{}
}

func (this *dynamicContext) appendSql(sqlStr string)  {
	this.sqlStr += sqlStr + " "
}

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
	if ok := exprProcess(this.test, ctx.params); ok {
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

	ctx.appendSql(this.open + " ")

	if ok {
		val := reflect.ValueOf(collection)

		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return false
		}

		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			keys := make([]string, 0)
			params := make(map[string]interface{})
			switch v.Kind() {
			case reflect.Struct :

			case reflect.Map:
				m := v.Interface().(map[string]interface{})
				for k, v := range m{
					key := this.item + "." + k
					keys = append(keys, key)
					params[key] = v
				}
			default:
				keys = append(keys, this.item)
				params[this.item] = v.Interface()
			}

			tempCtx := &dynamicContext{
				params: params,
			}

			this.sqlNode.build(tempCtx)

			this.tokenHandler(tempCtx, i)

			if i != 0 {
				ctx.appendSql(this.separator + " ")
			}

			ctx.appendSql(tempCtx.sqlStr)

			for _, k := range keys{
				delete(tempCtx.params, k)
			}

			for k, v := range tempCtx.params {
				ctx.params[k] = v
			}
		}

	}

	ctx.appendSql(this.close + " ")
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
			if string([]byte{sqlStr[i-1], sqlStr[i]}) == "#{"{
				start = i
			}
		}

		if start != 0 && i < len(sqlStr) -1 && sqlStr[i + 1] == '}'{
			finalSqlStr += sqlStr[:start+1]
			sqlStr = sqlStr[i+2:]

			var re = regexp.MustCompile("^\\s*" + this.item + "\\s*")
			itemPrefix := listItemPrefix + this.item + fmt.Sprintf("%d", index)
			s := re.ReplaceAllString(itemStr, itemPrefix)
			s = strings.Trim(s, " ")
			if strings.Contains(s, itemPrefix) {
				itemKey := strings.Trim(itemStr, " ")
				ctx.params[s] = ctx.params[itemKey]
			}

			finalSqlStr += s + "}"

			i = 0
			start = 0
			itemStr = ""
		}
	}

	finalSqlStr += sqlStr
	ctx.sqlStr = finalSqlStr
}


// set node

// trim node

// where node

// choose node
