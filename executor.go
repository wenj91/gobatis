package gobatis

import (
	"context"
	"errors"
	"fmt"
)

type executor struct {
	gb *gbBase
}

func (exec *executor) updateContext(ctx context.Context, ms *mappedStmt, params map[string]interface{}) (lastInsertId int64, affected int64, err error) {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return 0, 0, err
	}

	if conf.dbConf.ShowSQL {
		LOG.Info("SQL:%s\nParamMappings:%s\nParams:%v", boundSql.sqlStr, boundSql.paramMappings, paramArr)
	}

	stmt, err := exec.gb.db.PrepareContext(ctx, boundSql.sqlStr)
	if nil != err {
		return 0, 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, paramArr...)
	if nil != err {
		return 0, 0, err
	}

	lastInsertId, err = result.LastInsertId()
	if nil != err {
		return 0, 0, err
	}
	affected, err = result.RowsAffected()
	if nil != err {
		return 0, 0, err
	}

	return lastInsertId, affected, nil
}

func (exec *executor) update(ms *mappedStmt, params map[string]interface{}) (lastInsertId int64, affected int64, err error) {
	return exec.updateContext(context.Background(), ms, params)
}

func (exec *executor) queryContext(ctx context.Context, ms *mappedStmt, params map[string]interface{}, res interface{}, rowBound ...*rowBounds) (int64, error) {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return 0, err
	}

	sqlStr := boundSql.sqlStr

	if conf.dbConf.ShowSQL {
		LOG.Info("SQL:%s\nParamMappings:%s\nParams:%v", sqlStr, boundSql.paramMappings, paramArr)
	}

	count := int64(0)
	if len(rowBound) > 0 {
		countSql := "SELECT COUNT(1) cnt FROM (" + sqlStr + ") AS t"
		rows, err := exec.gb.db.QueryContext(ctx, countSql, paramArr...)
		if nil != err {
			return 0, err
		}

		resProc, err := rowsToMaps(rows)
		if nil != err {
			return 0, err
		}

		c, err := valToInt64(resProc[0].(map[string]interface{})["cnt"])
		if nil != err {
			return 0, err
		}

		count = c

		if count <= 0 {
			return 0, nil
		}

		sqlStr += " LIMIT " + fmt.Sprint(rowBound[0].offset) + "," + fmt.Sprint(rowBound[0].limit)
	}

	rows, err := exec.gb.db.QueryContext(ctx, sqlStr, paramArr...)
	if nil != err {
		return 0, err
	}
	defer rows.Close()

	resProc, ok := resSetProcMap[ms.resultType]
	if !ok {
		return 0, errors.New("No exec result type proc, result type:" + string(ms.resultType))
	}

	// func(rows *sql.Rows, res interface{}) error
	err = resProc(rows, res)
	if nil != err {
		return 0, err
	}

	return count, nil
}

func (exec *executor) query(ms *mappedStmt, params map[string]interface{}, res interface{}, rowBound ...*rowBounds) (int64, error) {
	return exec.queryContext(context.Background(), ms, params, res, rowBound...)
}

func paramProc(ms *mappedStmt, params map[string]interface{}) (boundSql *boundSql, paramArr []interface{}, err error) {
	boundSql = ms.sqlSource.getBoundSql(params)
	if nil == boundSql {
		err = errors.New("get boundSql err: boundSql == nil")
		return
	}

	paramArr = make([]interface{}, 0)
	for i := 0; i < len(boundSql.paramMappings); i++ {
		paramName := boundSql.paramMappings[i]
		param, ok := boundSql.extParams[paramName]
		if !ok {
			err = errors.New("param:" + paramName + " not exists")
			return
		}

		paramArr = append(paramArr, param)
	}

	return
}
