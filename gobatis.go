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
type Executor interface {
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

var showSql = false

type Config struct {
	Db          *sql.DB
	MapperPaths []string
}

func NewGoBatis(ctx context.Context, conf *Config) (*Gobatis, error) {
	if nil == conf.Db {
		panic("")
	}
	mapper := loadingMapper(conf.MapperPaths...)

	gb := &Gobatis{
		db:      conf.Db,
		mappers: mapper,
	}
	gb.ctxStd, gb.cancel = context.WithCancel(ctx)
	return gb, nil
}

type Runner struct {
	executor Executor
	dbType   DbType
	mappers  *mapper
	tx       *sql.Tx
}

// Gobatis
type Gobatis struct {
	ctxStd  context.Context
	cancel  context.CancelFunc
	mappers *mapper
	db      *sql.DB
}

// Begin Tx
//
// ps：
//  Tx, err := this.Begin()
func (g *Gobatis) Begin() (*Runner, error) {
	return g.BeginTx(g.ctxStd, nil)
}

// Begin Tx with ctx & opts
//
// ps：
//  Tx, err := this.BeginTx(ctx, ops)
func (g *Gobatis) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Runner, error) {
	tx, err := g.db.BeginTx(ctx, opts)
	if nil != err {
		return nil, err
	}
	return &Runner{
		mappers: g.mappers,
		tx:      tx,
	}, nil
}

// Close db
//
// ps：
//  err := this.Close()
func (g *Gobatis) Close() error {
	g.cancel()
	return g.db.Close()
}

// Commit Tx
//
// ps：
//  err := Tx.Commit()
func (r *Runner) Commit() error {
	if nil == r.tx {
		return errors.New("tx no running")
	}
	return r.tx.Commit()
}

// Rollback Tx
//
// ps：
//  err := Tx.Rollback()
func (r *Runner) Rollback() error {
	return r.tx.Rollback()
}

// reference from https://github.com/yinshuwei/osm/blob/master/osm.go end
func (r *Runner) Select(stmt string, param interface{}) func(res interface{}) error {
	ms := r.mappers.getMappedStmt(stmt)
	if nil == ms {
		return func(res interface{}) error {
			return errors.New("mapped statement not found:" + stmt)
		}
	}
	ms.dbType = r.dbType

	params := paramProcess(param)

	return func(res interface{}) error {
		executor := &executor{r}
		return executor.query(ms, params, res)
	}
}

// insert(Executor string, param interface{})
func (r *Runner) Insert(stmt string, param interface{}) (int64, int64, error) {
	ms := r.mappers.getMappedStmt(stmt)
	if nil == ms {
		return 0, 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = r.dbType

	params := paramProcess(param)

	executor := &executor{r}

	lastInsertId, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, 0, err
	}

	return lastInsertId, affected, nil
}

// update(Executor string, param interface{})
func (r *Runner) Update(stmt string, param interface{}) (int64, error) {
	ms := r.mappers.getMappedStmt(stmt)
	if nil == ms {
		return 0, errors.New("Mapped statement not found:" + stmt)
	}
	ms.dbType = r.dbType
	params := paramProcess(param)

	executor := &executor{r}

	_, affected, err := executor.update(ms, params)
	if nil != err {
		return 0, err
	}

	return affected, nil
}

// delete(Executor string, param interface{})
func (r *Runner) Delete(stmt string, param interface{}) (int64, error) {
	return r.Update(stmt, param)
}
