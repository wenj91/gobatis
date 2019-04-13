package gobatis

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

func createSqlNode(elems ...element) []iSqlNode {
	res := make([]iSqlNode, 0)
	if len(elems) == 0 {
		res = append(res, &textSqlNode{""})
		return res
	}

	if len(elems) == 1 {
		elem := elems[0]
		if elem.ElementType == eleTpText {
			res = append(res, &textSqlNode{
				content: elem.Val.(string),
			})

			return res
		}

		n := elem.Val.(node)
		if n.Name == "if" || n.Name == "when" {
			sqlNodes := createSqlNode(n.Elements...)
			ifn := &ifSqlNode{
				test: n.Attrs["test"].Value,
			}

			ifn.sqlNode = sqlNodes[0]
			if len(sqlNodes) > 1 {
				ifn.sqlNode = &mixedSqlNode{
					sqlNodes: sqlNodes,
				}
			}

			res = append(res, ifn)
			return res
		}

		if n.Name == "choose" {
			sqlNodes := createSqlNode(n.Elements...)
			csNode := &chooseNode{
				sqlNodes: sqlNodes,
			}
			res = append(res, csNode)
			return res
		}

		if n.Name == "otherwise" {
			sqlNodes := createSqlNode(n.Elements...)
			owNode := &mixedSqlNode{
				sqlNodes: sqlNodes,
			}
			res = append(res, owNode)
			return res
		}

		if n.Name == "foreach" {
			open := ""
			openAttr, ok := n.Attrs["open"]
			if ok {
				open = openAttr.Value
			}

			closeStr := ""
			closeAttr, ok := n.Attrs["close"]
			if ok {
				closeStr = closeAttr.Value
			}

			separator := ""
			separatorAttr, ok := n.Attrs["separator"]
			if ok {
				separator = separatorAttr.Value
			}

			itemAttr, ok := n.Attrs["item"]
			if !ok {
				log.Fatalln("No attr:`item` for tag:", n.Name)
			}
			item := itemAttr.Value

			index := ""
			indexAttr, ok := n.Attrs["index"]
			if ok {
				index = indexAttr.Value
			}

			collectionAttr, ok := n.Attrs["collection"]
			if !ok {
				log.Fatalln("No attr:`collection` for tag:", n.Name)
			}
			collection := collectionAttr.Value

			sqlNodes := createSqlNode(n.Elements...)

			fn := &foreachSqlNode{
				open:       open,
				close:      closeStr,
				separator:  separator,
				item:       item,
				index:      index,
				collection: collection,
			}

			fn.sqlNode = sqlNodes[0]
			if len(sqlNodes) > 1 {
				fn.sqlNode = &mixedSqlNode{
					sqlNodes: sqlNodes,
				}
			}

			res = append(res, fn)
			return res
		}

		if n.Name == "set" {
			sqlNodes := createSqlNode(n.Elements...)
			setN := &setSqlNode{
				sqlNodes: sqlNodes,
			}

			res = append(res, setN)
			return res
		}

		if n.Name == "trim" {
			sqlNodes := createSqlNode(n.Elements...)

			prefix := ""
			prefixAttr, ok := n.Attrs["prefix"]
			if ok {
				prefix = prefixAttr.Value
			}

			preOv := ""
			preOvAttr, ok := n.Attrs["prefixOverrides"]
			if ok {
				preOv = preOvAttr.Value
			}

			suffOv := ""
			suffOvAttr, ok := n.Attrs["suffixOverrides"]
			if ok {
				suffOv = suffOvAttr.Value
			}

			suffix := ""
			suffixAttr, ok := n.Attrs["suffix"]
			if ok {
				suffix = suffixAttr.Value
			}
			trimN := &trimSqlNode{
				sqlNodes:        sqlNodes,
				prefix:          prefix,
				prefixOverrides: preOv,
				suffixOverrides: suffOv,
				suffix:          suffix,
			}

			res = append(res, trimN)
			return res
		}

		if n.Name == "where" {
			sqlNodes := createSqlNode(n.Elements...)
			whereN := &whereSqlNode{
				sqlNodes: sqlNodes,
			}

			res = append(res, whereN)
			return res
		}

		log.Fatalln("The tag:", n.Name, "not support, current version only support tag:<if> | <foreach>")
	}

	for _, elem := range elems {
		sqlNode := createSqlNode(elem)
		res = append(res, sqlNode...)
	}

	return res
}

type dynamicContext struct {
	sqlStr string
	params map[string]interface{}
}

func newDynamicContext(params map[string]interface{}) *dynamicContext {
	return &dynamicContext{
		params: params,
	}
}

func (d *dynamicContext) appendSql(sqlStr string) {
	d.sqlStr += sqlStr + " "
}

func (d *dynamicContext) toSql() string {
	return d.sqlStr
}

// [ref](http://www.mybatis.org/mybatis-3/dynamic-sql.html)
type iSqlNode interface {
	build(ctx *dynamicContext) bool
}

// mixed node
type mixedSqlNode struct {
	sqlNodes []iSqlNode
}

func (m *mixedSqlNode) build(ctx *dynamicContext) bool {
	for i := 0; i < len(m.sqlNodes); i++ {
		sqlNode := m.sqlNodes[i]
		sqlNode.build(ctx)
	}

	return true
}

// if node
type ifSqlNode struct {
	test    string
	sqlNode iSqlNode
}

func (i *ifSqlNode) build(ctx *dynamicContext) bool {
	if ok := eval(i.test, ctx.params); ok {
		i.sqlNode.build(ctx)
		return true
	}

	return false
}

// text node
type textSqlNode struct {
	content string
}

func (t *textSqlNode) build(ctx *dynamicContext) bool {
	ctx.appendSql(t.content)
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

func (f *foreachSqlNode) build(ctx *dynamicContext) bool {
	collection, ok := ctx.params[f.collection]
	if !ok {
		log.Println("No collection for foreach tag:", f.collection)
		return false
	}

	ctx.appendSql(f.open)

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
			m := f.structToMap(v.Interface())
			for k, v := range m {
				key := f.item + "." + k
				keys = append(keys, key)
				params[key] = v
			}
		case reflect.Map:
			m := v.Interface().(map[string]interface{})
			for k, v := range m {
				key := f.item + "." + k
				keys = append(keys, key)
				params[key] = v
			}
		default:
			keys = append(keys, f.item)
			params[f.item] = v.Interface()
		}

		params[f.item] = v.Interface()

		tempCtx := &dynamicContext{
			params: params,
		}

		f.sqlNode.build(tempCtx)
		f.tokenHandler(tempCtx, i)

		if i != 0 {
			ctx.appendSql(f.separator)
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
	ctx.appendSql(f.close)

	return true
}

func (f *foreachSqlNode) tokenHandler(ctx *dynamicContext, index int) {
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

			var re = regexp.MustCompile("^\\s*" + f.item + "\\s*")
			itemPrefix := listItemPrefix + f.item + fmt.Sprintf("%d", index)
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

func (f *foreachSqlNode) structToMap(s interface{}) map[string]interface{} {
	return structToMap(s)
}

// set node
type setSqlNode struct {
	sqlNodes []iSqlNode
}

func (s *setSqlNode) build(ctx *dynamicContext) bool {

	sqlStr := ""
	for _, sqlNode := range s.sqlNodes {
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

func (t *trimSqlNode) build(ctx *dynamicContext) bool {
	tempCtx := &dynamicContext{
		params: ctx.params,
	}

	for _, sqlNode := range t.sqlNodes {
		if tempCtx.sqlStr != "" {
			tempCtx.sqlStr += " "
		}
		sqlNode.build(tempCtx)
	}

	if tempCtx.sqlStr != "" {
		sqlStr := strings.TrimSpace(tempCtx.sqlStr)

		preOv := strings.TrimSpace(t.prefixOverrides)
		if preOv != "" {
			sqlStr = strings.TrimPrefix(sqlStr, preOv+" ")
		}

		suffOv := strings.TrimSpace(t.suffixOverrides)
		if suffOv != "" {
			sqlStr = strings.TrimSuffix(sqlStr, suffOv+" ")
		}

		pre := strings.TrimSpace(t.prefix)
		if pre != "" {
			sqlStr = pre + " " + sqlStr
		}

		suff := strings.TrimSpace(t.suffix)
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

func (w *whereSqlNode) build(ctx *dynamicContext) bool {
	tempCtx := &dynamicContext{
		params: ctx.params,
	}

	for _, sqlNode := range w.sqlNodes {
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
type chooseNode struct {
	sqlNodes  []iSqlNode
	otherwise iSqlNode
}

func (c *chooseNode) build(ctx *dynamicContext) bool {
	for _, n := range c.sqlNodes {
		if n.build(ctx) {
			return true
		}
	}
	if nil != c.otherwise {
		c.otherwise.build(ctx)
		return true
	}
	return false
}
