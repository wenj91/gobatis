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

func NewGobatis() *gobatis {
	if nil == conf {
		log.Fatalln("Db config no init, please invoke gobatis.ConfInit() to init db config!")
		panic(errors.New("Db config no init, please invoke gobatis.ConfInit() to init db config!"))
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

	gb := &gobatis{
		GbBase{
			db:     db,
			dbType: DbType(conf.dbConf.DB.DriverName),
			config: conf,
		},
	}

	return gb
}

type GbBase struct {
	db     dbRunner
	dbType DbType
	config *config
}

// gobatis
type gobatis struct {
	GbBase
}

// tx
type tx struct {
	GbBase
}

// Begin Tx
//
// ps：
//  tx, err := this.Begin()
func (this *GbBase) Begin() (*tx, error) {
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

	t := &tx{
		GbBase{
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
//  tx, err := this.BeginTx(ctx, ops)
func (this *GbBase) BeginTx(ctx context.Context, opts *sql.TxOptions) (*tx, error) {
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

	t := &tx{
		GbBase{
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
func (this *GbBase) Close() error {
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
//  err := tx.Commit()
func (this *tx) Commit() error {
	if nil == this.db {
		return errors.New("tx no running")
	}

	sqlTx, ok := this.db.(*sql.Tx)
	if !ok {
		return errors.New("tx no running")

	}

	return sqlTx.Commit()
}

// Rollback tx
//
// ps：
//  err := tx.Rollback()
func (this *tx) Rollback() error {
	if nil == this.db {
		return errors.New("tx no running")
	}

	sqlTx, ok := this.db.(*sql.Tx)
	if !ok {
		return errors.New("tx no running")
	}

	return sqlTx.Rollback()
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go end

func (this *GbBase) Select(stmt string, param interface{}) func(res interface{}) error {
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
func (this *GbBase) Insert(stmt string, param interface{}) (int64, error) {
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
func (this *GbBase) Update(stmt string, param interface{}) (int64, error) {
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
func (this *GbBase) Delete(stmt string, param interface{}) (int64, error) {
	return this.Update(stmt, param)
}
