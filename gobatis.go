package gobatis

import (
	"context"
	"database/sql"
	"errors"
	"github.com/wenj91/gobatis/logger"
	"github.com/wenj91/gobatis/m"
	"github.com/wenj91/gobatis/sb"
	param2 "github.com/wenj91/gobatis/uti/param"
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
	Wrapper(model m.Model) Wrapper
	// Select 查询数据
	Select(stmt string, param interface{}) func(res interface{}) error
	// SelectContext 查询数据with context
	SelectContext(ctx context.Context, stmt string, param interface{}) func(res interface{}) error
	// Insert 插入数据
	Insert(stmt string, param interface{}) (lastInsertId int64, affected int64, err error)
	// InsertContext 插入数据with context
	InsertContext(ctx context.Context, stmt string, param interface{}) (lastInsertId int64, affected int64, err error)
	// Update 更新数据
	Update(stmt string, param interface{}) (affected int64, err error)
	// UpdateContext 更新数据with context
	UpdateContext(ctx context.Context, stmt string, param interface{}) (affected int64, err error)
	// Delete 刪除数据
	Delete(stmt string, param interface{}) (affected int64, err error)
	// DeleteContext 刪除数据with context
	DeleteContext(ctx context.Context, stmt string, param interface{}) (affected int64, err error)
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go start
type dbRunner interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func Get(datasource string) *DB {
	if nil == conf {
		panic(errors.New("DB config no init, please invoke DB.ConfInit() to init db config!"))
	}

	if nil == db {
		panic(errors.New("DB init err, db == nil!"))
	}

	ds, ok := db[datasource]
	if !ok {
		panic(errors.New("Datasource:" + datasource + " not exists!"))
	}

	dbType := ds.dbType
	if dbType != DBTypeMySQL {
		panic(errors.New("No support to this driver name!"))
	}

	gb := &DB{
		gbBase{
			db:     ds.db,
			dbType: ds.dbType,
			config: conf,
		},
	}

	return gb
}

func SetLogger(log logger.Logger) {
	logger.LOG = log
}

type gbBase struct {
	db     dbRunner
	dbType DBType
	config *config
}

// DB
type DB struct {
	gbBase
}

var _ GoBatis = &DB{}

// TX
type TX struct {
	gbBase
}

var _ GoBatis = &TX{}

type Wrapper struct {
	gb GoBatis
	wp sb.Wrapper
}

// Begin TX
//
// ps：
//  TX, err := this.Begin()
func (d *DB) Begin() (*TX, error) {
	if nil == d.db {
		return nil, errors.New("db no opened")
	}

	sqlDB, ok := d.db.(*sql.DB)
	if !ok {
		return nil, errors.New("db no opened")
	}

	db, err := sqlDB.Begin()
	if nil != err {
		return nil, err
	}

	t := &TX{
		gbBase{
			dbType: d.dbType,
			config: d.config,
			db:     db,
		},
	}
	return t, nil
}

// Begin TX with ctx & opts
//
// ps：
//  TX, err := this.BeginTx(ctx, ops)
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TX, error) {
	if nil == d.db {
		return nil, errors.New("db no opened")
	}

	sqlDb, ok := d.db.(*sql.DB)
	if !ok {
		return nil, errors.New("db no opened")
	}

	db, err := sqlDb.BeginTx(ctx, opts)
	if nil != err {
		return nil, err
	}

	t := &TX{
		gbBase{
			dbType: d.dbType,
			config: d.config,
			db:     db,
		},
	}
	return t, nil
}

// Transaction tx
func (d *DB) Transaction(fn func(tx *TX) error) error {

	tx, err := d.Begin()
	if nil != err {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if nil != err {
				logger.LOG.Error("tx rollback err:#v", err)
			}
		}
	}()

	err = fn(tx)
	if nil != err {
		return err
	}

	err = tx.Commit()
	if nil != err {
		return err
	}

	return nil
}

// Transaction tx
func (d *DB) TransactionTX(ctx context.Context, opts *sql.TxOptions, fn func(tx *TX) error) error {

	tx, err := d.BeginTx(ctx, opts)
	if nil != err {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if nil != err {
				logger.LOG.Error("tx rollback err:#v", err)
			}
		}
	}()

	err = fn(tx)
	if nil != err {
		return err
	}

	err = tx.Commit()
	if nil != err {
		return err
	}

	return nil
}

func (g *gbBase) Wrapper(model m.Model) Wrapper {
	return Wrapper{
		gb: g,
		wp: sb.Model(model),
	}
}

// Close db
//
// ps：
//  err := this.Close()
func (g *gbBase) Close() error {
	if nil == g.db {
		return errors.New("db no opened")
	}

	sqlDb, ok := g.db.(*sql.DB)
	if !ok {
		return errors.New("db no opened")
	}

	err := sqlDb.Close()
	g.db = nil
	return err
}

// Commit TX
//
// ps：
//  err := TX.Commit()
func (t *TX) Commit() error {
	if nil == t.db {
		return errors.New("TX no running")
	}

	sqlTx, ok := t.db.(*sql.Tx)
	if !ok {
		return errors.New("TX no running")

	}

	return sqlTx.Commit()
}

// Rollback TX
//
// ps：
//  err := TX.Rollback()
func (t *TX) Rollback() error {
	if nil == t.db {
		return errors.New("TX no running")
	}

	sqlTx, ok := t.db.(*sql.Tx)
	if !ok {
		return errors.New("TX no running")
	}

	return sqlTx.Rollback()
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go end
func (g *gbBase) Select(stmt string, param interface{}) func(res interface{}) error {
	ms := g.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return func(res interface{}) error {
			return errors.New("Mapped statement not found:" + stmt)
		}
	}
	ms.dbType = g.dbType

	params := param2.Process(param)

	return func(res interface{}) error {
		executor := &executor{
			gb: g,
		}
		err := executor.query(ms, params, res)
		return err
	}
}

func (g *gbBase) SelectContext(ctx context.Context, stmt string, param interface{}) func(res interface{}) error {
	ms := g.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return func(res interface{}) error {
			return errors.New("Mapped statement not found:" + stmt)
		}
	}
	ms.dbType = g.dbType

	params := param2.Process(param)

	return func(res interface{}) error {
		executor := &executor{
			gb: g,
		}
		err := executor.queryContext(ctx, ms, params, res)
		return err
	}
}

// insert(stmt string, param interface{})
func (g *gbBase) Insert(stmt string, param interface{}) (int64, int64, error) {
	ms := g.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = g.dbType

	params := param2.Process(param)

	executor := &executor{
		gb: g,
	}

	lastInsertId, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, 0, err
	}

	return lastInsertId, affected, nil
}

func (g *gbBase) InsertContext(ctx context.Context, stmt string, param interface{}) (int64, int64, error) {
	ms := g.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = g.dbType

	params := param2.Process(param)

	executor := &executor{
		gb: g,
	}

	lastInsertId, affected, err := executor.updateContext(ctx, ms, params)
	if nil != err {
		return 0, 0, err
	}

	return lastInsertId, affected, nil
}

// update(stmt string, param interface{})
func (g *gbBase) Update(stmt string, param interface{}) (int64, error) {
	ms := g.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = g.dbType

	params := param2.Process(param)

	executor := &executor{
		gb: g,
	}

	_, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, err
	}

	return affected, nil
}

func (g *gbBase) UpdateContext(ctx context.Context, stmt string, param interface{}) (int64, error) {
	ms := g.config.mapperConf.getMappedStmt(stmt)
	if nil == ms {
		return 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = g.dbType

	params := param2.Process(param)

	executor := &executor{
		gb: g,
	}

	_, affected, err := executor.updateContext(ctx, ms, params)
	if nil != err {
		return 0, err
	}

	return affected, nil
}

// delete(stmt string, param interface{})
func (g *gbBase) Delete(stmt string, param interface{}) (int64, error) {
	return g.Update(stmt, param)
}

func (g *gbBase) DeleteContext(ctx context.Context, stmt string, param interface{}) (int64, error) {
	return g.UpdateContext(ctx, stmt, param)
}
