package mapperstmt

import (
	"github.com/wenj91/gobatis/constants"
	"github.com/wenj91/gobatis/xmlparser"
	"io"
)

type SqlNode struct {
	xmlparser.Node
	IsDynamic  bool
	ResultType string
}

type StmtMapper struct {
	Namespace   string
	SelectStmts map[string]SqlNode
	InsertStmts map[string]SqlNode
	UpdateStmts map[string]SqlNode
	DeleteStmts map[string]SqlNode
}

func GetStmtMapper(r io.Reader) StmtMapper {
	root := xmlparser.Parse(r)

	res := StmtMapper{}

	rootNode := root
	for _, val := range rootNode.Attr {
		if val.Name.Local == "namespace" {
			res.Namespace = val.Value
			break
		}
	}

	selectStmts := make(map[string]SqlNode)
	insertStmts := make(map[string]SqlNode)
	updateStmts := make(map[string]SqlNode)
	deleteStmts := make(map[string]SqlNode)

	eles := rootNode.Elements
	for _, val := range eles {
		if val.ElementType == constants.ELE_TP_NODE {
			eleNode := val.Val.(xmlparser.Node)
			sqlNode := SqlNode{
				Node:      eleNode,
				IsDynamic: false,
			}
			//judge isDynamic
			for _, val := range eleNode.Elements {
				if val.ElementType == constants.ELE_TP_NODE {
					sqlNode.IsDynamic = true
					break
				}
			}
			for key, val := range eleNode.Attr {
				if key == "resultType" {
					sqlNode.ResultType = val.Value
				}
			}
			switch eleNode.Name {
			case "insert":
				check(insertStmts, eleNode.Id)
				insertStmts[eleNode.Id] = sqlNode
			case "delete":
				check(deleteStmts, eleNode.Id)
				deleteStmts[eleNode.Id] = sqlNode
			case "update":
				check(updateStmts, eleNode.Id)
				updateStmts[eleNode.Id] = sqlNode
			case "select":
				check(selectStmts, eleNode.Id)
				selectStmts[eleNode.Id] = sqlNode
			}
		}
	}

	res.InsertStmts = insertStmts
	res.DeleteStmts = deleteStmts
	res.UpdateStmts = updateStmts
	res.SelectStmts = selectStmts

	return res
}

func check(sqlNode map[string]SqlNode, id string) {
	res := sqlNode[id]
	if res.Name != "" {
		panic("init err, find multi id:" + id)
	}
}
