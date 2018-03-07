package gobatis

import (
	"context"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	"database/sql"
	"encoding/xml"

	"github.com/wenj91/gobatis/constants"
	"github.com/wenj91/gobatis/process/resultprocess"
	"github.com/wenj91/gobatis/tools/datautil"
	"github.com/wenj91/gobatis/tools/regutil"

	_ "github.com/go-sql-driver/mysql"
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

// select insert update delete tag
type stmt struct {
	Id         string `xml:"id,attr"`
	ResultType string `xml:"resultType,attr"`
	SQL        string `xml:",chardata"`
}

// mapper tag
type mapper struct {
	Namespace   string `xml:"namespace,attr"`
	SelectStmts []stmt `xml:"select"`
	InsertStmts []stmt `xml:"insert"`
	UpdateStmts []stmt `xml:"update"`
	DeleteStmts []stmt `xml:"delete"`
}

// sql mapper
type sqlMapper struct {
	Id         string
	SQL        string
	ResultType string
}

type goBatis struct {
	sqlMappers map[string]sqlMapper
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
		sqlMappers: make(map[string]sqlMapper),
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
	sqlStr, sqlParams, resultType, err := sqlProcess(gobatis.sqlMappers[id], params...)
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
		stmt, err := gobatis.prepare(sqlStr)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		defer stmt.Close()

		rows, err := stmt.Query(sqlParams...)
		if nil != err {
			log.Println(err)
			return 0, err
		}

		if len(results) == 1 {
			switch resultType {
			case constants.RESULT_TYPE_MAP:
				return resultprocess.MapProcess(rows, results[0], sqlParams)
			case constants.RESULT_TYPE_MAPS:
				return resultprocess.MapsProcess(rows, results[0], sqlParams)
			case constants.RESULT_TYPE_SLICE:
				return resultprocess.SliceProcess(rows, results[0], sqlParams)
			case constants.RESULT_TYPE_SLICES:
				return resultprocess.SlicesProcess(rows, results[0], sqlParams)
			case constants.RESULT_TYPE_STRUCT:
				return resultprocess.StructProcess(rows, results[0], sqlParams)
			case constants.RESULT_TYPE_STRUCTS:
				return resultprocess.StructsProcess(rows, results[0], sqlParams)
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
	sqlStr, sqlParams, _, err := sqlProcess(gobatis.sqlMappers[id], params...)
	if nil != err {
		log.Println(err)
		return 0, 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	stmt, err := gobatis.prepare(sqlStr)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(sqlParams...)
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

	sqlStr, sqlParams, _, err := sqlProcess(gobatis.sqlMappers[id], params...)
	if nil != err {
		log.Println(err)
		return 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, err
	}

	stmt, err := gobatis.prepare(sqlStr)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(sqlParams...)
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
	sqlStr, sqlParams, _, err := sqlProcess(gobatis.sqlMappers[id], params...)
	if nil != err {
		log.Println(err)
		return 0, err
	}

	if err != nil {
		log.Println(err)
		return 0, err
	}

	stmt, err := gobatis.prepare(sqlStr)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(sqlParams...)
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

// sql参数提取并生成prepare sql语句和参数列表
// @result
//		string: sqlStr生成prepare sql语句
//		interface{}: sqlParams生成prepare参数列表
//		string: resultType返回类型 目前支持：Map, Maps, Struct, Structs, Slice, Slices这几种形式返回
func sqlProcess(sqlMapper sqlMapper, params ...interface{}) (sqlStr string, sqlParams []interface{}, resultType string, err error) {

	if sqlMapper.Id == "" {
		err = errors.New("no this id in mapper file")
		return
	}

	//如果sqlTemplate中存在参数, 则提取参数
	paramNames := []string{}
	if strings.Contains(sqlMapper.SQL, "#{") {
		paramNames = regutil.SharpParamNamesFind(sqlMapper.SQL)
	}

	//简化sql语句
	sqlStr = strings.Replace(sqlMapper.SQL, "\r", " ", -1)
	sqlStr = strings.Replace(sqlStr, "\n", " ", -1)
	sqlStr = strings.Replace(sqlStr, "\t", " ", -1)
	sqlStr = strings.Trim(sqlStr, " ")

	//转化sql语句
	for i := 0; i < len(paramNames); i++ {
		sqlStr, err = regutil.SharpParamMatchReplace(sqlStr, paramNames[i])
		if nil != err {
			return
		}
	}

	resultType = sqlMapper.ResultType

	var param interface{}
	paramsSize := len(params)
	if paramsSize > 0 {
		if paramsSize == 1 {
			param = params[0]
		} else {
			param = params
		}

		paramVal := reflect.ValueOf(param)
		kind := paramVal.Kind()
		switch {
		case kind == reflect.Array || kind == reflect.Slice:
			for i := 0; i < paramVal.Len() && i < len(paramNames); i++ {
				itemVal := paramVal.Index(i)
				sqlParams = append(sqlParams, itemVal.Interface())
			}
		case kind == reflect.Map:
			paramMap := param.(map[string]interface{})
			for i := 0; i < len(paramNames); i++ {
				item := paramMap[paramNames[i]]
				if nil == item {
					err = errors.New("params must not be nil")
					return
				}
				sqlParams = append(sqlParams, item)
			}
		case kind == reflect.Struct:
			paramVal := reflect.ValueOf(param)
			if paramVal.Kind() == reflect.Ptr {
				err = errors.New("struct params must not be ptr")
				return
			}
			for i := 0; i < len(paramNames); i++ {
				item := datautil.FieldToParams(param, paramNames[i])
				if nil == item {
					err = errors.New("no this params:" + paramNames[i])
					return
				}
				sqlParams = append(sqlParams, item)
			}

		case kind == reflect.Bool ||
			kind == reflect.Int ||
			kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32 ||
			kind == reflect.Int64 ||
			kind == reflect.Uint ||
			kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32 ||
			kind == reflect.Uint64 ||
			kind == reflect.Uintptr ||
			kind == reflect.Float32 ||
			kind == reflect.Float64 ||
			kind == reflect.Complex64 ||
			kind == reflect.Complex128 ||
			kind == reflect.String:
			sqlParams = append(sqlParams, param)
		}
	}

	log.Println("sql:", sqlStr)
	log.Println("sqlParams:", sqlParams)

	return
}

// 处理mapper xml文件，将其转化成
func (gobatis *goBatis) sqlMapperProcess(name string) {
	file, err := os.Open(name)
	if nil != err {
		log.Println(err)
		panic(err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	mapper := mapper{}
	decoder.Decode(&mapper)

	// 解析生成insert stmt模板
	for _, insertStmt := range mapper.InsertStmts {
		sm := stmtProcess(mapper, insertStmt)
		gobatis.sqlMappers[sm.Id] = sm
	}

	// 解析生成delete stmt模板
	for _, deleteStmt := range mapper.DeleteStmts {
		sm := stmtProcess(mapper, deleteStmt)
		gobatis.sqlMappers[sm.Id] = sm
	}

	// 解析生成update stmt模板
	for _, updateStmt := range mapper.UpdateStmts {
		sm := stmtProcess(mapper, updateStmt)
		gobatis.sqlMappers[sm.Id] = sm
	}

	// 解析生成select stmt模板
	for _, selectStmt := range mapper.SelectStmts {
		sm := stmtProcess(mapper, selectStmt)
		gobatis.sqlMappers[sm.Id] = sm
	}
}

// 处理stmt，生成sqlMapper
func stmtProcess(mapper mapper, stmt stmt) (sm sqlMapper) {
	sqlMapper := sqlMapper{}

	id := stmt.Id
	if mapper.Namespace != "" {
		id = mapper.Namespace + "#" + id
	}
	resultType := stmt.ResultType
	sqlStr := stmt.SQL

	sqlMapper.Id = id
	sqlMapper.SQL = sqlStr
	sqlMapper.ResultType = resultType

	log.Println("sql mapper init: -->", id, resultType, sqlStr)

	return sqlMapper
}
