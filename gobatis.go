package gobatis

import (
	"context"
	"database/sql"
	"fmt"
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
	mapperConfig mapperConfig
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
