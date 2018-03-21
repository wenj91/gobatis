package sqlsource

import (
	"github.com/wenj91/boundsql"
	"github.com/wenj91/mapperstmt"
	"errors"
	"strings"
	"log"
)

type SqlSource struct {
	sqlNode mapperstmt.SqlNode
}

func NewSqlSource(sqlNode mapperstmt.SqlNode) SqlSource {
	if sqlNode.Id == "" {
		panic("can not find this mapper id")
	}
	return SqlSource{
		sqlNode: sqlNode,
	}
}

//方案:
//    如果是静态节点,则参数可以支持是数组形式
//    如果是动态节点,则参数只能是map,dto或者entity
func (sqlSource *SqlSource) GetBoundSql(params ...interface{}) (res boundsql.BoundSql, err error) {
	isDynamic := sqlSource.sqlNode.IsDynamic

	paramSize := len(params)
	var param interface{}
	if paramSize == 1 {
		param = params[0]
	} else {
		param = params
	}

	if isDynamic {
		// 动态sql不可能没有参数的
		if paramSize == 0 {
			err = errors.New("param size must be more than 0")
			return
		}
		//动态sql生成
		res, err = boundsql.DynamicGetBoundSql(sqlSource.sqlNode, param)
	} else {
		log.Println("param size is:", paramSize)
		//静态sql生成
		res, err = boundsql.StaticGetBoundSql(sqlSource.sqlNode, param)
	}

	//简化sql语句
	sqlStr := res.Sql
	sqlStr = strings.Replace(sqlStr, "\r", " ", -1)
	sqlStr = strings.Replace(sqlStr, "\n", " ", -1)
	sqlStr = strings.Replace(sqlStr, "\t", " ", -1)
	sqlStr = strings.Trim(sqlStr, " ")

	res.Sql = sqlStr
	res.ResultType = sqlSource.sqlNode.ResultType

	log.Println(res)

	return
}
