package gobatis

import (
	"context"
	"errors"
	"github.com/wenj91/gobatis/logger"
	"github.com/wenj91/gobatis/m"
)

type executor struct {
	gb *gbBase
}

func (exec *executor) wrapperUpdateContext(ctx context.Context, sqlStr string, paramArr []interface{}) (lastInsertId int64, affected int64, err error) {
	stmt, err := exec.gb.db.PrepareContext(ctx, sqlStr)
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

func (exec *executor) updateContext(ctx context.Context, ms *mappedStmt, params map[string]interface{}) (lastInsertId int64, affected int64, err error) {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return 0, 0, err
	}

	if conf.dbConf.ShowSQL {
		logger.LOG.Info("SQL:%s\nParamMappings:%s\nParams:%v", boundSql.sqlStr, boundSql.paramMappings, paramArr)
	}

	return exec.wrapperUpdateContext(ctx, boundSql.sqlStr, paramArr)
}

func (exec *executor) wrapperUpdate(sqlStr string, paramArr []interface{}) (lastInsertId int64, affected int64, err error) {
	stmt, err := exec.gb.db.Prepare(sqlStr)
	if nil != err {
		return 0, 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(paramArr...)
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
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return 0, 0, err
	}

	if conf.dbConf.ShowSQL {
		logger.LOG.Info("SQL:%s\nParamMappings:%s\nParams:%v", boundSql.sqlStr, boundSql.paramMappings, paramArr)
	}

	return exec.wrapperUpdate(boundSql.sqlStr, paramArr)
}

func (exec *executor) wrapperQueryContext(ctx context.Context, sqlStr string, rt m.ResultType, paramArr []interface{}, res interface{}) error {
	rows, err := exec.gb.db.QueryContext(ctx, sqlStr, paramArr...)
	if nil != err {
		return err
	}
	defer rows.Close()

	resProc, ok := resSetProcMap[rt]
	if !ok {
		return errors.New("No exec result type proc, result type:" + string(rt))
	}

	// func(rows *sql.Rows, res interface{}) error
	err = resProc(rows, res)
	if nil != err {
		return err
	}

	return nil
}

func (exec *executor) queryContext(ctx context.Context, ms *mappedStmt, params map[string]interface{}, res interface{}) error {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return err
	}

	if conf.dbConf.ShowSQL {
		logger.LOG.Info("SQL:%s\nParamMappings:%s\nParams:%v", boundSql.sqlStr, boundSql.paramMappings, paramArr)
	}

	return exec.wrapperQueryContext(ctx, boundSql.sqlStr, ms.resultType, paramArr, res)
}

func (exec *executor) wrapperQuery(sqlStr string, rt m.ResultType, paramArr []interface{}, res interface{}) error {
	rows, err := exec.gb.db.Query(sqlStr, paramArr...)
	if nil != err {
		return err
	}
	defer rows.Close()

	resProc, ok := resSetProcMap[rt]
	if !ok {
		return errors.New("No exec result type proc, result type:" + string(rt))
	}

	// func(rows *sql.Rows, res interface{}) error
	err = resProc(rows, res)
	if nil != err {
		return err
	}

	return nil
}

func (exec *executor) query(ms *mappedStmt, params map[string]interface{}, res interface{}) error {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return err
	}

	if conf.dbConf.ShowSQL {
		logger.LOG.Info("SQL:%s\nParamMappings:%s\nParams:%v", boundSql.sqlStr, boundSql.paramMappings, paramArr)
	}

	return exec.wrapperQuery(boundSql.sqlStr, ms.resultType, paramArr, res)
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
