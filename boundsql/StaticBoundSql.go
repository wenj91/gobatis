package boundsql

import (
	"github.com/wenj91/gobatis/mapperstmt"
)

func StaticGetBoundSql(sqlNode mapperstmt.SqlNode, param interface{}) (boundSql BoundSql, err error) {

	eles := sqlNode.Elements
	sqlStr := ""
	for i:=0; i<len(eles); i++ {
		sqlStr += (eles[i].Val.(string))
	}

	boundSql, err = BaseGetBoundSql(sqlStr, param)
	return
}

