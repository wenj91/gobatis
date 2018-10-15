package gobatis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ResultType string

const (
	// result set is a map: map[string]interface{}
	resultTypeMap ResultType = "Map"
	// result set is a slice, item is map: []map[string]interface{}
	resultTypeMaps ResultType = "Maps"
	// result set is a struct
	resultTypeStruct ResultType = "Struct"
	// result set is a slice, item is struct
	resultTypeStructs ResultType = "Structs"
	// result set is a value slice, []interface{}
	resultTypeSlice ResultType = "Slice"
	// result set is a value slice, item is value slice, []interface{}
	resultTypeSlices ResultType = "Slices"
	// result set is single value
	resultTypeValue ResultType = "Value"
)

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go start
type dbRunner interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DbType string

const (
	DbTypeMysql    DbType = "mysql"
	DbTypePostgres DbType = "postgres"
)

type gbBase struct {
	db           dbRunner
	mapperConfig *mapperConfig
	dbType       DbType
}

// gobatis
type gobatis struct {
	gbBase
}

// tx
type tx struct {
	gbBase
}

// Begin Tx
//
// ps：
//  tx, err := this.Begin()
func (this *gbBase) Begin() (*tx, error) {
	t := &tx{}
	t.dbType = this.dbType

	if this.db == nil {
		err := fmt.Errorf("db no opened")
		return nil, err
	}

	sqlDb, ok := this.db.(*sql.DB)
	if !ok {
		err := fmt.Errorf("db no opened")
		return nil, err
	}

	db, err := sqlDb.Begin()
	if nil != err {
		return nil, err
	}

	t.db = db
	return t, nil
}

// Begin Tx with ctx & opts
//
// ps：
//  tx, err := this.BeginTx(ctx, ops)
func (this *gbBase) BeginTx(ctx context.Context, opts *sql.TxOptions) (*tx, error) {
	t := &tx{}
	t.dbType = this.dbType

	if this.db == nil {
		err := fmt.Errorf("db no opened")
		return nil, err
	}

	sqlDb, ok := this.db.(*sql.DB)
	if !ok {
		err := fmt.Errorf("db no opened")
		return nil, err
	}

	db, err := sqlDb.BeginTx(ctx, opts)
	if nil != err {
		return nil, err
	}

	t.db = db
	return t, nil
}

// Close db
//
// ps：
//  err := this.Close()
func (this *gbBase) Close() error {
	if this.db == nil {
		err := fmt.Errorf("db no opened")
		return err
	}

	sqlDb, ok := this.db.(*sql.DB)
	if !ok {
		err := fmt.Errorf("db no opened")
		return err
	}

	err := sqlDb.Close()
	this.db = nil
	return err
}

// Commit Tx
//
// ps：
//  err := tx.Commit()
func (this *tx) Commit() error {
	if this.db == nil {
		err := fmt.Errorf("tx no runing")
		return err
	}

	sqlTx, ok := this.db.(*sql.Tx)
	if !ok {
		err := fmt.Errorf("tx no runing")
		return err

	}

	return sqlTx.Commit()
}

// Rollback tx
//
// ps：
//  err := tx.Rollback()
func (this *tx) Rollback() error {
	if this.db == nil {
		err := fmt.Errorf("tx no runing")
		return err
	}

	sqlTx, ok := this.db.(*sql.Tx)
	if !ok {
		err := fmt.Errorf("tx no runing")
		return err
	}

	return sqlTx.Rollback()
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go end

func (this *gbBase) Select(stmt string, param interface{}) func(res interface{}) error {
	ms := this.mapperConfig.getMappedStmt(stmt)
	if nil != ms {
		return func(res interface{}) error {
			return errors.New("Mapped statement not found")
		}
	}

	ms.dbType = this.dbType

	params := paramProcess(param)

	return func(res interface{}) error {
		executor := &executor{
			gb: this,
		}
		err := executor.query(ms, params, res)
		return err
	}
}

// insert(stmt string, param interface{})
// update(stmt string, param interface{})
// delete(stmt string, param interface{})
