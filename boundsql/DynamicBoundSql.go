package boundsql

import (
	"github.com/wenj91/gobatis/constants"
	"github.com/wenj91/gobatis/mapperstmt"
	"github.com/wenj91/gobatis/process/tagprocess"
	"github.com/wenj91/gobatis/xmlparser"
)

func DynamicGetBoundSql(sqlNode mapperstmt.SqlNode, param interface{}) (boundsql BoundSql, err error) {

	sqlStr := ""
	sqlParams := make([]interface{}, 0)

	//遍历节点处理节点
	eles := sqlNode.Elements
	for i := 0; i < len(eles); i++ {
		if eles[i].ElementType == constants.ELE_TP_STRING {//字符数据处理
			sqlStr += (eles[i].Val.(string))
			//如果是单纯字符串,直接使用基本处理就可以
			boundSqlTmp, errTmp := BaseGetBoundSql(sqlStr, param)
			if nil != errTmp {
				err = errTmp
				return
			}
			sqlStr = boundSqlTmp.Sql
			sqlParams = append(sqlParams, boundSqlTmp.ParameterMappings...)
		} else {//节点数据处理
			//node process
			node := eles[i].Val.(xmlparser.Node)
			nodeType := node.Name
			switch nodeType {
			case "foreach":
				subSqlStr, subSqlParams, errTmp := tagprocess.ForeachProcess(node, param)
				if nil != errTmp {
					err = errTmp
					return
				}
				sqlStr += subSqlStr
				sqlParams = append(sqlParams, subSqlParams...)
			case "if":
				//todo:

			}
		}
	}

	boundsql.Sql = sqlStr
	boundsql.ParameterMappings = sqlParams

	return
}
