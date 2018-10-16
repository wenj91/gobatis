package gobatis

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
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

func NewGobatis() *Db {
	if nil == conf {
		log.Fatalln("Db config no init, please invoke Db.ConfInit() to init db config!")
		panic(errors.New("Db config no init, please invoke Db.ConfInit() to init db config!"))
	}

	db, err := sql.Open(conf.dbConf.DB.DriverName, conf.dbConf.DB.DataSourceName)
	if nil != err {
		log.Println(err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Println(err)
		panic(err)
	}

	if conf.dbConf.DB.MaxLifeTime == 0 {
		db.SetConnMaxLifetime(120 * time.Second)
	} else {
		db.SetConnMaxLifetime(time.Duration(conf.dbConf.DB.MaxLifeTime) * time.Second)
	}

	if conf.dbConf.DB.MaxOpenConns == 0 {
		db.SetMaxOpenConns(10)
	} else {
		db.SetMaxOpenConns(conf.dbConf.DB.MaxOpenConns)
	}

	if conf.dbConf.DB.MaxOpenConns == 0 {
		db.SetMaxIdleConns(5)
	} else {
		db.SetMaxIdleConns(conf.dbConf.DB.MaxIdleConns)
	}

	gb := &Db{
		gbBase{
			db:     db,
			dbType: DbType(conf.dbConf.DB.DriverName),
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
func (this *gbBase) Begin() (*Tx, error) {
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
func (this *gbBase) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
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
