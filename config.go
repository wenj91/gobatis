package gobatis

import (
	"database/sql"
	"github.com/wenj91/gobatis/logger"
	"os"
	"strings"
	"time"
)

type config struct {
	dbConf     *DBConfig
	mapperConf *mapperConfig
}

var conf *config
var db map[string]*GoBatisDB

func Init(option IOption) {
	configInit(option.ToDBConf())
}

func configInit(dbConf *DBConfig) {
	if nil == dbConf {
		panic("Build db config err: dbConf == nil")
	}

	if len(dbConf.DB) <= 0 && dbConf.db == nil {
		panic("No datasource config")
	}

	mapperConf := &mapperConfig{
		mappedStmts: make(map[string]*node),
		mappedSql:   make(map[string]*node),
	}

	for _, item := range dbConf.Mappers {
		f, err := os.Open(item)
		if nil != err {
			panic("Open mapper config: " + item + " err:" + err.Error())
		}

		logger.LOG.Info("mapper config:%s %s", item, "init...")
		mc := buildMapperConfig(f)
		for k, ms := range mc.mappedStmts {
			mapperConf.put(k, ms)
		}

		// sql tag cache
		for k, ms := range mc.mappedSql {
			mapperConf.putSql(k, ms)
		}
	}

	conf = &config{
		dbConf:     dbConf,
		mapperConf: mapperConf,
	}

	// init db
	dbInit(dbConf)
}

func dbInit(dbConf *DBConfig) {
	db = make(map[string]*GoBatisDB)
	if len(dbConf.DB) <= 0 && dbConf.db == nil {
		panic("No config for datasource")
	}

	for _, item := range dbConf.DB {
		if item.DataSource == "" {
			panic("DB config err: datasource must not be nil")
		}

		item.DataSource = strings.TrimSpace(item.DataSource)

		_, ok := db[item.DataSource]
		if ok {
			panic("DB config datasource name repeat:" + item.DataSource)
		}

		if item.DriverName == "" {
			panic("DB config err: driverName must not be nil")
		}

		if item.DataSourceName == "" {
			panic("DB config err: dataSourceName must not be nil")
		}

		dbConn, err := sql.Open(item.DriverName, item.DataSourceName)
		if nil != err {
			panic(err)
		}

		if err := dbConn.Ping(); err != nil {
			panic(err)
		}

		if item.MaxLifeTime == 0 {
			dbConn.SetConnMaxLifetime(120 * time.Second)
		} else {
			dbConn.SetConnMaxLifetime(time.Duration(item.MaxLifeTime) * time.Second)
		}

		if item.MaxOpenConns == 0 {
			dbConn.SetMaxOpenConns(10)
		} else {
			dbConn.SetMaxOpenConns(item.MaxOpenConns)
		}

		if item.MaxOpenConns == 0 {
			dbConn.SetMaxIdleConns(5)
		} else {
			dbConn.SetMaxIdleConns(item.MaxIdleConns)
		}

		d := NewGoBatisDB(DBType(item.DriverName), dbConn)
		db[item.DataSource] = d
	}

	if dbConf.db != nil {
		for k, v := range dbConf.db {
			_, ok := db[k]
			if ok {
				panic("DB config datasource name repeat:" + k)
			}
			db[k] = v
		}
	}
}
