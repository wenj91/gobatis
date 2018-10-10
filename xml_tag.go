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

	if start != 0{
		log.Println("WARN: token not close")
	}

	finalSqlStr += sqlStr
	ctx.sqlStr = finalSqlStr
}

func (this *foreachSqlNode) structToMap(s interface{}) map[string]interface{} {
	objVal := reflect.ValueOf(s)
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	res := make(map[string]interface{})
	objType := objVal.Type()
	for i:=0; i<objVal.NumField(); i++{
		fieldVal := objVal.Field(i)
		if fieldVal.CanInterface() {
			field := objType.Field(i)
			res[field.Name]=fieldVal.Interface()
		}
	}

	return res
}

// set node

// trim node

// where node

// choose node
