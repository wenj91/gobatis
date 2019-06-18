package gobatis

import (
	"context"
	"database/sql"
	"errors"
)

type ResultType string

const (
	resultTypeMap     ResultType = "map"     // result set is a map: map[string]interface{}
	resultTypeMaps    ResultType = "maps"    // result set is a slice, item is map: []map[string]interface{}
	resultTypeStruct  ResultType = "struct"  // result set is a struct
	resultTypeStructs ResultType = "structs" // result set is a slice, item is struct
	resultTypeSlice   ResultType = "slice"   // result set is a value slice, []interface{}
	resultTypeSlices  ResultType = "slices"  // result set is a value slice, item is value slice, []interface{}
	resultTypeArray   ResultType = "array"   //
	resultTypeArrays  ResultType = "arrays"  // result set is a value slice, item is value slice, []interface{}
	resultTypeValue   ResultType = "value"   // result set is single value
)

type GoBatis interface {
	Select(stmt string, param interface{}) func(res interface{}) error
	Insert(stmt string, param interface{}) (int64, int64, error)
	Update(stmt string, param interface{}) (int64, error)
	Delete(stmt string, param interface{}) (int64, error)
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go start
type dbRunner interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DbType string

const (
	dbTypeMySQL    DbType = "mysql"
	dbTypePostgres DbType = "postgres"
)

func NewGoBatis(datasource string) *Db {
	if nil == conf {
		panic(errors.New("Db config no init, please invoke Db.ConfInit() to init db config!"))
	}

	if nil == db {
		panic(errors.New("Db init err, db == nil!"))
	}

	ds, ok := db[datasource]
	if !ok {
		panic(errors.New("Datasource:" + datasource + "not exists!"))
	}

	dbType := conf.dbConf.getDataSourceByName(datasource).DriverName
	if dbType != "mysql" {
		panic(errors.New("No support to this driver name!"))
	}

	gb := &Db{
		gbBase{
			db:     ds,
			dbType: DbType(dbType),
			config: conf,
		},
	}

	return gb
}

type gbBase struct {
	db     dbRunner
	dbType DbType
	config *config
}

// Db
type Db struct {
	gbBase
}

// Tx
type Tx struct {
	gbBase
}

// Begin Tx
//
// ps：
//  Tx, err := this.Begin()
func (this *Db) Begin() (*Tx, error) {
	if nil == this.db {
		return nil, errors.New("db no opened")
	}

	sqlDb, ok := this.db.(*sql.DB)
	if !ok {
		return nil, errors.New("db no opened")
	}

	db, err := sqlDb.Begin()
	if nil != err {
		return nil, err
	}

	t := &Tx{
		gbBase{
			dbType: this.dbType,
			config: this.config,
			db:     db,
		},
	}
	return t, nil
}

// Begin Tx with ctx & opts
//
// ps：
//  Tx, err := this.BeginTx(ctx, ops)
func (this *Db) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if nil == this.db {
		return nil, errors.New("db no opened")
	}

	sqlDb, ok := this.db.(*sql.DB)
	if !ok {
		return nil, errors.New("db no opened")
	}

	db, err := sqlDb.BeginTx(ctx, opts)
	if nil != err {
		return nil, err
	}

	t := &Tx{
		gbBase{
			dbType: this.dbType,
			config: this.config,
			db:     db,
		},
	}
	return t, nil
}

// Close db
//
// ps：
//  err := this.Close()
func (this *gbBase) Close() error {
	if nil == this.db {
		return errors.New("db no opened")
	}

	sqlDb, ok := this.db.(*sql.DB)
	if !ok {
		return errors.New("db no opened")
	}

	err := sqlDb.Close()
	this.db = nil
	return err
}

// Commit Tx
//
// ps：
//  err := Tx.Commit()
func (this *Tx) Commit() error {
	if nil == this.db {
		return errors.New("Tx no running")
	}

	sqlTx, ok := this.db.(*sql.Tx)
	if !ok {
		return errors.New("Tx no running")

	}

	return sqlTx.Commit()
}

// Rollback Tx
//
// ps：
//  err := Tx.Rollback()
func (this *Tx) Rollback() error {
	if nil == this.db {
		return errors.New("Tx no running")
	}

	sqlTx, ok := this.db.(*sql.Tx)
	if !ok {
		return errors.New("Tx no running")
	}

	return sqlTx.Rollback()
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go end

func (this *gbBase) Select(stmt string, param interface{}) func(res interface{}) error {
	ms := this.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return func(res interface{}) error {
			return errors.New("Mapped statement not found:" + stmt)
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
func (this *gbBase) Insert(stmt string, param interface{}) (int64, int64, error) {
	ms := this.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = this.dbType

	params := paramProcess(param)

	executor := &executor{
		gb: this,
	}

	lastInsertId, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, 0, err
	}

	return lastInsertId, affected, nil
}

// update(stmt string, param interface{})
func (this *gbBase) Update(stmt string, param interface{}) (int64, error) {
	ms := this.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = this.dbType

	params := paramProcess(param)

	executor := &executor{
		gb: this,
	}

	_, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, err
	}

	return affected, nil
}

// delete(stmt string, param interface{})
func (this *gbBase) Delete(stmt string, param interface{}) (int64, error) {
	return this.Update(stmt, param)
}
