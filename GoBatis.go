package gobatis

import (
	"github.com/wenj91/gobatis/mapperstmt"
	"github.com/wenj91/gobatis/sqlsource"
	"context"
	"database/sql"
	"errors"
	"github.com/wenj91/gobatis/constants"
	"github.com/wenj91/gobatis/process/resultprocess"
	"log"
	"os"
)

//next todo:
//0 todo: 单元测试编写
//1 todo: 批量插入修改<for>标签实现
//2 todo: 动态sql生成<if>标签实现
//3 todo: ${xxx}解析实现
//4 todo: 结果集映射<resultMap>标签实现
//5 todo: 公共查询字段<sql>标签实现
//6 todo: 一级缓存实现
//7 todo: 二级缓存实现
//8 todo: 完善文档

type goBatis struct {
	sqlMappers map[string]mapperstmt.SqlNode
	db         *sql.DB
	isBegin    bool
	tx         *sql.Tx
}

func NewGoBatis(driverName string, url string, mappers []string) goBatis {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	db, err := sql.Open(driverName, url)
	if nil != err {
		log.Println(err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Println(err)
		panic(err)
	}

	gobatis := goBatis{
		sqlMappers: make(map[string]mapperstmt.SqlNode),
		db:         db,
	}

	// mapper init
	for _, xmlPath := range mappers {
		gobatis.sqlMapperProcess(xmlPath)
	}

	return gobatis
}

// tx begin
func (gobatis *goBatis) Begin() error {
	tx, err := gobatis.db.Begin()
	gobatis.isBegin = true
	gobatis.tx = tx
	return err
}

// tx begin with ctx&opts
func (gobatis *goBatis) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	tx, err := gobatis.db.BeginTx(ctx, opts)
	gobatis.isBegin = true
	gobatis.tx = tx
	return err
}

// tx commit
func (gobatis *goBatis) Commit() error {
	err := gobatis.tx.Commit()
	if nil == err {
		gobatis.isBegin = false
		gobatis.tx = nil
	}
	return err
}

// tx rollback
func (gobatis *goBatis) Rollback() error {
	if gobatis.isBegin {
		err := gobatis.tx.Rollback()
		if nil == err {
			gobatis.isBegin = false
			gobatis.tx = nil
		}
		return err
	}

	return errors.New("no tx in this op")
}

// 根据是否开启事务选择Prepare操作
func (gobatis *goBatis) prepare(sqlStr string) (*sql.Stmt, error) {
	if gobatis.isBegin {
		return gobatis.tx.Prepare(sqlStr)
	} else {
		return gobatis.db.Prepare(sqlStr)
	}
}

// @method:
// 		select
// @params:
//		id: xml 映射文件id
//		params: 传入参数, 可以是: array || slice || map || struct || base type(int unint...)
// @results:
//		return type: map || maps || slice || slices || struct || structs
//		results: 返回结果
//		int: 返回数据条数
//		error: 错误信息
func (gobatis *goBatis) Select(id string, params ...interface{}) func(results ...interface{}) (int, error) {
	ss := sqlsource.NewSqlSource(gobatis.sqlMappers[id])
	boundSql, err := ss.GetBoundSql(params...)
	if nil != err {
		log.Println(err)
		return func(results ...interface{}) (int, error) {
			return 0, err
		}
	}

	if err != nil {
		log.Println(err)
		return func(results ...interface{}) (int, error) {
			return 0, err
		}
	}

	callback := func(results ...interface{}) (int, error) {
		stmt, err := gobatis.prepare(boundSql.Sql)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		defer stmt.Close()

		rows, err := stmt.Query(boundSql.ParameterMappings...)
		if nil != err {
			log.Println(err)
			return 0, err
		}

		if len(results) == 1 {
			switch boundSql.ResultType {
			case constants.RESULT_TYPE_MAP:
				return resultprocess.MapProcess(rows, results[0], boundSql.ParameterMappings)
			case constants.RESULT_TYPE_MAPS:
				return resultprocess.MapsProcess(rows, results[0], boundSql.ParameterMappings)
			case constants.RESULT_TYPE_SLICE:
				return resultprocess.SliceProcess(rows, results[0], boundSql.ParameterMappings)
			case constants.RESULT_TYPE_SLICES:
				return resultprocess.SlicesProcess(rows, results[0], boundSql.ParameterMappings)
			case constants.RESULT_TYPE_STRUCT:
				return resultprocess.StructProcess(rows, results[0], boundSql.ParameterMappings)
			case constants.RESULT_TYPE_STRUCTS:
				return resultprocess.StructsProcess(rows, results[0], boundSql.ParameterMappings)
			default:
				return 0, errors.New("no this result type define")
			}
		}

		return 0, errors.New("result size must be 1")
	}

	return callback
}

// 保存操作
// @method:
// 		insert
// @params:
//		id: xml 映射文件id
//		params: 传入参数, 可以是: array || slice || map || struct || base type(int unint...)
// @results:
//		int: lastInsertId, 返回插入Id
//		int: affectedCount, 返回插入数据条数
//		error: 错误信息
func (gobatis *goBatis) Insert(id string, params ...interface{}) (int, int, error) {
	ss := sqlsource.NewSqlSource(gobatis.sqlMappers[id])
	boundSql, err := ss.GetBoundSql(params)
	if nil != err {
		log.Println(err)
		return 0, 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	stmt, err := gobatis.prepare(boundSql.Sql)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(boundSql.ParameterMappings...)
	if nil != err {
		log.Println(err)
		return 0, 0, err
	}

	insertId, err := result.LastInsertId()
	if nil != err {
		return 0, 0, err
	}
	affected, err := result.RowsAffected()
	if nil != err {
		return 0, 0, err
	}
	return int(insertId), int(affected), err
}

//func (gobatis *goBatis) InsertBatch(id string, params ...interface{}) (int, error) {
//
//	//todo:
//
//	return 0, nil
//}

// 删除操作
// @method:
// 		insert
// @params:
//		id: xml 映射文件id
//		params: 传入参数, 可以是: array || slice || map || struct || base type(int unint...)
// @results:
//		int: lastInsertId, 返回插入Id
//		int: affectedCount, 返回插入数据条数
//		error: 错误信息
func (gobatis *goBatis) Delete(id string, params ...interface{}) (int, error) {

	ss := sqlsource.NewSqlSource(gobatis.sqlMappers[id])
	boundSql, err := ss.GetBoundSql(params)
	if nil != err {
		log.Println(err)
		return 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, err
	}

	stmt, err := gobatis.prepare(boundSql.Sql)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(boundSql.ParameterMappings...)
	if nil != err {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if nil != err {
		return 0, err
	}

	return int(affected), nil
}

func (gobatis *goBatis) Update(id string, params ...interface{}) (int, error) {
	ss := sqlsource.NewSqlSource(gobatis.sqlMappers[id])
	boundSql, err := ss.GetBoundSql(params)
	if nil != err {
		log.Println(err)
		return 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, err
	}

	stmt, err := gobatis.prepare(boundSql.Sql)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(boundSql.ParameterMappings...)
	if nil != err {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if nil != err {
		return 0, err
	}

	return int(affected), nil
}

//func (gobatis *goBatis) UpdateBatch(id string, params ...interface{}) (int, error) {
//	//todo:
//
//	return 0, nil
//}

// 处理mapper xml文件，将其转化成
func (gobatis *goBatis) sqlMapperProcess(name string) {
	file, err := os.Open(name)
	if nil != err {
		log.Println(err)
		panic(err)
	}
	defer file.Close()

	mapper := mapperstmt.GetStmtMapper(file)

	// 解析生成insert stmt模板
	for _, insertStmt := range mapper.InsertStmts {
		id := getMapperId(mapper, insertStmt)
		gobatis.sqlMappers[id] = insertStmt
	}

	// 解析生成delete stmt模板
	for _, deleteStmt := range mapper.DeleteStmts {
		id := getMapperId(mapper, deleteStmt)
		gobatis.sqlMappers[id] = deleteStmt
	}

	// 解析生成update stmt模板
	for _, updateStmt := range mapper.UpdateStmts {
		id := getMapperId(mapper, updateStmt)
		gobatis.sqlMappers[id] = updateStmt
	}

	// 解析生成select stmt模板
	for _, selectStmt := range mapper.SelectStmts {
		id := getMapperId(mapper, selectStmt)
		gobatis.sqlMappers[id] = selectStmt
	}
}

func getMapperId(mapper mapperstmt.StmtMapper, node mapperstmt.SqlNode) string {
	id := ""
	if mapper.Namespace != "" {
		id += (mapper.Namespace + ".")
	}
	id += node.Id

	return id
}
