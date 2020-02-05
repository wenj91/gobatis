package gobatis

import (
	"io"
	"strings"

	"gopkg.in/yaml.v2"
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
				LOG.Error("No attr:`item` for tag:%s", n.Name)
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
				LOG.Error("No attr:`collection` for tag:%s", n.Name)
				panic("No attr:`collection` for tag:" + n.Name)
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

		LOG.Error("The tag:" + n.Name + "not support, current version only support tag:<if> | <foreach>")
		panic("The tag:" + n.Name + "not support, current version only support tag:<if> | <foreach>")
	}

	for _, elem := range elems {
		sqlNode := createSqlNode(elem)
		res = append(res, sqlNode...)
	}

	return res
}

func buildMapperConfig(r io.Reader) *mapperConfig {
	rootNode := parse(r)

	conf := &mapperConfig{
		mappedStmts: make(map[string]*node),
		cache:       make(map[string]*mappedStmt),
	}

	if rootNode.Name != "mapper" {
		LOG.Error("Mapper xml must start with `mapper` tag, please check your xml mapperConfig!")
		panic("Mapper xml must start with `mapper` tag, please check your xml mapperConfig!")
	}

	namespace := ""
	if val, ok := rootNode.Attrs["namespace"]; ok {
		nStr := strings.TrimSpace(val.Value)
		if nStr != "" {
			nStr += "."
		}
		namespace = nStr
	}

	for _, elem := range rootNode.Elements {
		if elem.ElementType == eleTpNode {
			childNode := elem.Val.(node)
			switch childNode.Name {
			case "select", "update", "insert", "delete":
				if childNode.Id == "" {
					LOG.Error("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
					panic("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
				}

				fid := namespace + childNode.Id
				if ok := conf.put(fid, &childNode); !ok {
					LOG.Error("Repeat id for:" + fid + "Please check your xml mapperConfig!")
					panic("Repeat id for:" + fid + "Please check your xml mapperConfig!")
				}

			case "sql":
				if childNode.Id == "" {
					LOG.Error("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
					panic("No id for:" + childNode.Name + "Id must be not null, please check your xml mapperConfig!")
				}

				fid := namespace + childNode.Id
				if ok := conf.put(fid, &childNode); !ok {
					LOG.Error("Repeat id for:" + fid + "Please check your xml mapperConfig!")
					panic("Repeat id for:" + fid + "Please check your xml mapperConfig!")
				}
			}
		}
	}

	return conf
}

func buildDbConfig(ymlStr string) *DBConfig {
	dbconf := &DBConfig{}
	err := yaml.Unmarshal([]byte(ymlStr), &dbconf)
	if err != nil {
		panic("error: " + err.Error())
	}

	return dbconf
}
