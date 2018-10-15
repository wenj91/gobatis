package gobatis

import (
	"errors"
	"fmt"
)

type executor struct {
	gb *gbBase
}


func (this *executor) query(ms *mappedStmt, params map[string]interface{}, res interface{}) error {
	boundSql := ms.sqlSource.getBoundSql(params)
	fmt.Println("SQL:", boundSql.sqlStr)
	fmt.Println("ParamMappings:", boundSql.paramMappings)

	paramArr := make([]interface{}, 0)
	for i := 0; i < len(boundSql.paramMappings); i++ {
		paramName := boundSql.paramMappings[i]
		param, ok := boundSql.extParams[paramName]
		if !ok {
			return errors.New("param:" + paramName + " not exists")
		}

		paramArr = append(paramArr, param)
	}

	fmt.Println("Params:", paramArr)

	rows, err := this.gb.db.Query(boundSql.sqlStr, paramArr...)
	if nil != err {
		return err
	}

	resProc, ok := resSetProcMap[ms.resultType]
	if !ok {
		return errors.New("No this result type proc, result type:" + string(ms.resultType))
	}

	// func(rows *sql.Rows, res interface{}) error
	err = resProc(rows, res)
	if nil != err {
		return err
	}

	return nil
}

