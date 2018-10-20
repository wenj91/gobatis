package gobatis

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type ResultType string

const (
	// result set is a map: map[string]interface{}
	resultTypeMap ResultType = "Map"
	resultTypeMapL ResultType = "map"
	// result set is a slice, item is map: []map[string]interface{}
	resultTypeMaps ResultType = "Maps"
	resultTypeMapsL ResultType = "maps"
	// result set is a struct
	resultTypeStruct ResultType = "Struct"
	resultTypeStructL ResultType = "struct"
	// result set is a slice, item is struct
	resultTypeStructs ResultType = "Structs"
	resultTypeStructsL ResultType = "structs"
	// result set is a value slice, []interface{}
	resultTypeSlice ResultType = "Slice"
	resultTypeSliceL ResultType = "slice"
	resultTypeArray ResultType = "array"
	// result set is a value slice, item is value slice, []interface{}
	resultTypeSlices ResultType = "Slices"
	resultTypeSlicesL ResultType = "slices"
	resultTypeArrays ResultType = "arrays"
	// result set is single value
	resultTypeValue ResultType = "Value"
	resultTypeValueL ResultType = "value"
)

type GoBatis interface {
	Select(stmt string, param interface{}) func(res interface{}) error
	Insert(stmt string, param interface{}) (int64, error)
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
		log.Fatalln("Db config no init, please invoke Db.ConfInit() to init db config!")
		panic(errors.New("Db config no init, please invoke Db.ConfInit() to init db config!"))
	}

	if nil == db {
		log.Fatalln("Db init err, db == nil!")
		panic(errors.New("Db init err, db == nil!"))
	}

	ds, ok := db[datasource]
	if !ok {
		log.Fatalln("Datasource:", datasource, "not exists!")
		panic(errors.New("Datasource:" + datasource + "not exists!"))
	}

	gb := &Db{
		gbBase{
			db:     ds,
			dbType: DbType(conf.dbConf.DB[datasource].DriverName),
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
func (this *gbBase) Insert(stmt string, param interface{}) (int64, error) {
	ms := this.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = this.dbType

	params := paramProcess(param)

	executor := &executor{
		gb: this,
	}

	lastInsertId, _, err := executor.update(ms, params)
	if nil != err {
		return 0, err
	}

	return lastInsertId, nil
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
