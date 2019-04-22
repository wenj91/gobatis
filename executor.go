package gobatis

import (
	"errors"
	"log"
)

type executor struct {
	gb *gbBase
}

func (exec *executor) update(ms *mappedStmt, params map[string]interface{}) (lastInsertId int64, affected int64, err error) {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return 0, 0, err
	}

	if conf.dbConf.ShowSql {
		log.Println("SQL:", boundSql.sqlStr)
		log.Println("ParamMappings:", boundSql.paramMappings)
		log.Println("Params:", paramArr)
	}

	stmt, err := exec.gb.db.Prepare(boundSql.sqlStr)
	if nil != err {
		return 0, 0, err
	}

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

func (exec *executor) query(ms *mappedStmt, params map[string]interface{}, res interface{}) error {
	boundSql, paramArr, err := paramProc(ms, params)
	if nil != err {
		return err
	}

	if conf.dbConf.ShowSql {
		log.Println("SQL:", boundSql.sqlStr)
		log.Println("ParamMappings:", boundSql.paramMappings)
		log.Println("Params:", paramArr)
	}

	rows, err := exec.gb.db.Query(boundSql.sqlStr, paramArr...)
	if nil != err {
		return err
	}

	resProc, ok := resSetProcMap[ms.resultType]
	if !ok {
		return errors.New("No exec result type proc, result type:" + string(ms.resultType))
	}

	// func(rows *sql.Rows, res interface{}) error
	err = resProc(rows, res)
	if nil != err {
		return err
	}

	return nil
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
