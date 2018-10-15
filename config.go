package gobatis

import (
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	dbConf     *dbConfig
	mapperConf *mapperConfig
}

var conf *config

func ConfInit(dbConfPath string)  {
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

	mapperConf := &mapperConfig{
		mappedStmts: make(map[string]*node),
	}
	for _, item := range dbConf.Mappers {
		f, err := os.Open(item)
		if nil != err {
			log.Fatalln("Open mapper config:", item, "err:", err)
			return
		}

		mc := buildMapperConfig(f)
		for k, ms := range mc.mappedStmts {
			mapperConf.put(k, ms)
		}
	}

	conf = &config{
		dbConf:     dbConf,
		mapperConf: mapperConf,
	}
}
