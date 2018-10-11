package gobatis

import "log"

func createSqlNode(elems ...element) iSqlNode {
	if len(elems) == 0 {
		return &textSqlNode{""}
	}

	if len(elems) == 1 {
		elem := elems[0]
		if elem.ElementType == eleTpText {
			return &textSqlNode{
				content: elem.Val.(string),
			}
		}

		n := elem.Val.(node)
		if n.Name == "if" {
			sqlNode := createSqlNode(n.Elements...)
			return &ifSqlNode{
				test:    n.Attrs["test"].Value,
				sqlNode: sqlNode,
			}
		}

		if n.Name == "for" {
			sqlNode := createSqlNode(n.Elements...)

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

			item := ""
			itemAttr, ok := n.Attrs["item"]
			if ok {
				item = itemAttr.Value
			}

			index := ""
			indexAttr, ok := n.Attrs["index"]
			if ok {
				index = indexAttr.Value
			}

			collection := ""
			collectionAttr, ok := n.Attrs["collection"]
			if ok {
				collection = collectionAttr.Value
			}

			return &foreachSqlNode{
				sqlNode:   sqlNode,
				open:      open,
				close:     closeStr,
				separator: separator,
				item:      item,
				index:     index,
				collection:collection,
			}
		}

		log.Fatalln("The tag:", n.Name, "not support, current version only support tag:<if> | <for>")
	}

	sns := make([]iSqlNode, 0)
	for _, elem := range elems {
		sqlNode := createSqlNode(elem)
		sns = append(sns, sqlNode)
	}

	return &mixedSqlNode{
		sqlNodes: sns,
	}
}
