package gobatis

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
	"time"
)

type TUser struct {
	Id       int64      `field:"id"`
	Name     string     `field:"name"`
	Password NullString `field:"password"`
}

func TestGoBatis(t *testing.T) {
	ConfInit("")
	if nil == conf {
		log.Println("db config == nil")
		return
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
		gbBase{
			db:           db,
			mapperConfig: conf.mapperConf,
		},
	}

	res, _ := gb.db.Query("select 1")
	cols, _ := res.Columns()
	fmt.Println(cols)

	//result := make(map[string]interface{})
	//result := make([]interface{}, 0)
	//var result interface{}
	//result := make([]TUser, 0)
	var result TUser
	err = gb.Select("Mapper.findById", map[string]interface{}{
		"id": 2,
	})(&result)

	fmt.Println("result:", result, "err:", err)

	u := &TUser{
		Name: "wenj1991",
		Password: NullString{
			String: "654321",
			Valid:  true,
		},
	}

	id, err := gb.Insert("Mapper.saveUser", u)
	fmt.Println("id:", id, "err:", err)

	uu := &TUser{
		Id:3,
		Name: "wenj1993",
		Password: NullString{
			String: "654321",
			Valid:  true,
		},
	}
	affected, err := gb.Update("Mapper.updateByStruct", uu)
	fmt.Println("affected:", affected, "err:", err)

	affected, err = gb.Delete("Mapper.deleteById", map[string]interface{}{
		"id":3,
	})
	fmt.Println("delete affected:", affected, "err:", err)
}
