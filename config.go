package gobatis

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type config struct {
	dbConf     *dbConfig
	mapperConf *mapperConfig
}

var conf *config
var db *sql.DB

func ConfInit(dbConfPath string)  {
	if nil != conf {
		log.Println("[WARN] Db config is already init, do not repeat init!")
		return
	}

	if dbConfPath == "" {
		dbConfPath = "db.yml"
	}
	f, err := os.Open(dbConfPath)
	if nil != err {
		log.Fatalln("Open db conf err:", err)
		return
	}

	r, err := ioutil.ReadAll(f)
	if nil != err {
		log.Fatalln("Read db conf err:", err)
		return
	}
	dbConf := buildDbConfig(string(r))
	if nil == dbConf {
		log.Fatalln("Build db config err: dbConf == nil")
		return
	}

	if dbConf.DB.DriverName == "" {
		log.Fatalln("Db config err: driverName must not be nil")
		return
	}

	if dbConf.DB.DataSourceName == "" {
		log.Fatalln("Db config err: dataSourceName must not be nil")
		return
	}

	mapperConf := &mapperConfig{
		mappedStmts: make(map[string]*node),
	}
	for _, item := range dbConf.Mappers {
		f, err := os.Open(item)
		if nil != err {
			log.Fatalln("Open mapper config:", item, "err:", err)
			return
		}

		log.Println("mapper config:", item, "init...")
		mc := buildMapperConfig(f)
		for k, ms := range mc.mappedStmts {
			mapperConf.put(k, ms)
		}
	}

	conf = &config{
		dbConf:     dbConf,
		mapperConf: mapperConf,
	}

	// init db
	dbInit()
}

func dbInit()  {
	dbConn, err := sql.Open(conf.dbConf.DB.DriverName, conf.dbConf.DB.DataSourceName)
	if nil != err {
		log.Println(err)
		panic(err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Println(err)
		panic(err)
	}

	if conf.dbConf.DB.MaxLifeTime == 0 {
		dbConn.SetConnMaxLifetime(120 * time.Second)
	} else {
		dbConn.SetConnMaxLifetime(time.Duration(conf.dbConf.DB.MaxLifeTime) * time.Second)
	}

	if conf.dbConf.DB.MaxOpenConns == 0 {
		dbConn.SetMaxOpenConns(10)
	} else {
		dbConn.SetMaxOpenConns(conf.dbConf.DB.MaxOpenConns)
	}

	if conf.dbConf.DB.MaxOpenConns == 0 {
		dbConn.SetMaxIdleConns(5)
	} else {
		dbConn.SetMaxIdleConns(conf.dbConf.DB.MaxIdleConns)
	}

	db = dbConn
}
