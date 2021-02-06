package gobatis

import (
	"github.com/wenj91/gobatis/logger"
	"sync"
)

type mapperConfig struct {
	mappedStmts map[string]*node
	mappedSql   map[string]*node
	cache       map[string]*mappedStmt
	mu          sync.Mutex
}

func (mc *mapperConfig) put(id string, n *node) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, ok := mc.mappedStmts[id]; ok {
		return false
	}

	mc.mappedStmts[id] = n
	return true
}

func (mc *mapperConfig) putSql(id string, n *node) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, ok := mc.mappedSql[id]; ok {
		return false
	}

	mc.mappedSql[id] = n
	return true
}

func (mc *mapperConfig) getXmlNode(id string) (rootNode *node, resultType string) {
	rootNode, ok := mc.mappedStmts[id]
	if !ok {
		panic("Can not find id:" + id + "mapped stmt")
	}

	resultType = ""
	if rootNode.Name == "select" {
		resultTypeAttr, ok := rootNode.Attrs["resultType"]
		if !ok {
			panic("Tag `<select>` must have resultType attr!")
		}

		resultType = resultTypeAttr.Value
	}

	return
}

func (mc *mapperConfig) getMappedStmt(id string) *mappedStmt {
	if nil == mc.cache {
		mc.cache = make(map[string]*mappedStmt)
	}

	if st, ok := mc.cache[id]; ok {
		return st
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	stmt := mc.buildSqlNode(id)
	mc.cache[id] = stmt

	return stmt
}

func (mc *mapperConfig) buildSqlNode(id string) *mappedStmt {
	rootNode, resultType := mc.getXmlNode(id)

	sn := mc.createSqlNode(rootNode.Elements...)

	ds := &dynamicSqlSource{}
	ds.sqlNode = sn[0]
	if len(sn) > 1 {
		ds.sqlNode = &mixedSqlNode{
			sqlNodes: sn,
		}
	}

	stmt := &mappedStmt{
		sqlSource:  ds,
		resultType: ResultType(resultType),
	}

	return stmt
}

func (mc *mapperConfig) createSqlNode(elems ...element) []iSqlNode {
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

		// include tag process
		if n.Name == "include" {
			id := n.getAttr("refid")
			id = n.Namespace + id
			ic, ok := mc.mappedSql[id]
			if !ok {
				logger.LOG.Error("No include sql for id:%s", id)
				panic("No include sql for id:" + id)
			}

			sqlNodes := mc.createSqlNode(ic.Elements...)
			res = append(res, sqlNodes...)
			return res
		}

		if n.Name == "if" || n.Name == "when" {
			sqlNodes := mc.createSqlNode(n.Elements...)
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
			sqlNodes := mc.createSqlNode(n.Elements...)
			csNode := &chooseNode{
				sqlNodes: sqlNodes,
			}
			res = append(res, csNode)
			return res
		}

		if n.Name == "otherwise" {
			sqlNodes := mc.createSqlNode(n.Elements...)
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
				logger.LOG.Error("No attr:`item` for tag:%s", n.Name)
				panic("No attr:`item` for tag:" + n.Name)
			}
			item := itemAttr.Value

			index := ""
			indexAttr, ok := n.Attrs["index"]
			if ok {
				index = indexAttr.Value
			}

			collectionAttr, ok := n.Attrs["collection"]
			if !ok {
				logger.LOG.Error("No attr:`collection` for tag:%s", n.Name)
				panic("No attr:`collection` for tag:" + n.Name)
			}
			collection := collectionAttr.Value

			sqlNodes := mc.createSqlNode(n.Elements...)

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
			sqlNodes := mc.createSqlNode(n.Elements...)
			setN := &setSqlNode{
				sqlNodes: sqlNodes,
			}

			res = append(res, setN)
			return res
		}

		if n.Name == "trim" {
			sqlNodes := mc.createSqlNode(n.Elements...)

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
			sqlNodes := mc.createSqlNode(n.Elements...)
			whereN := &whereSqlNode{
				sqlNodes: sqlNodes,
			}

			res = append(res, whereN)
			return res
		}

		logger.LOG.Error("The tag:" + n.Name + "not support, current version only support tag:<if> <when> <choose> <otherwise> <foreach> <set> <trim> <where>")
		panic("The tag:" + n.Name + "not support, current version only support tag:<if> <when> <choose> <otherwise> <foreach> <set> <trim> <where>")
	}

	for _, elem := range elems {
		sqlNode := mc.createSqlNode(elem)
		res = append(res, sqlNode...)
	}

	return res
}
