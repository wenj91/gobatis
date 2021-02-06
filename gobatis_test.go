package gobatis

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wenj91/gobatis/logger"
	"github.com/wenj91/gobatis/na"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type TUser struct {
	Id       int64         `field:"id"`
	Name     string        `field:"name"`
	Password na.NullString `field:"pwd"`
	Email    na.NullString `field:"email"`
	CrtTm    na.NullTime   `field:"crtTm"`
}

func TestGoBatis(t *testing.T) {
	Init(NewFileOption())
	if nil == conf {
		logger.LOG.Error("db config == nil")
		return
	}

	gb := Get("ds")

	//result := make(map[string]interface{})
	//result := make([]interface{}, 0)
	//var result interface{}
	//result := make([]TUser, 0)
	var result *TUser
	err := gb.Select("userMapper.findById", map[string]interface{}{
		"id": 2,
	})(&result)

	fmt.Println("result:", result, "err:", err)

	u := &TUser{
		Name: "wenj1991",
		Password: na.NullString{
			String: "654321",
			Valid:  true,
		},
	}

	id, a, err := gb.Insert("userMapper.saveUser", u)
	fmt.Println("id:", id, "affected:", a, "err:", err)

	uu := &TUser{
		Id:   1,
		Name: "wenj1993",
		Password: na.NullString{
			String: "654321",
			Valid:  true,
		},
	}

	// test set
	affected, err := gb.Update("userMapper.updateByCond", uu)
	fmt.Println("updateByCond:", affected, err)

	param := &TUser{
		Name: "wenj1993",
	}

	// test where
	res := make([]*TUser, 0)
	err = gb.Select("userMapper.queryStructsByCond", param)(&res)
	fmt.Println("queryStructsByCond", res, err)

	// test trim
	res2 := make([]*TUser, 0)
	err = gb.Select("userMapper.queryStructsByCond2", param)(&res2)
	fmt.Println("queryStructsByCond", res2, err)

	affected, err = gb.Delete("userMapper.deleteById", map[string]interface{}{
		"id": 3,
	})
	fmt.Println("delete affected:", affected, "err:", err)
}

func TestGoBatisWithDB(t *testing.T) {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	dbs := make(map[string]*GoBatisDB)
	dbs["ds"] = NewGoBatisDB(DBTypeMySQL, db)

	option := NewDBOption().
		DB(dbs).
		ShowSQL(true).
		Mappers([]string{"examples/mapper/userMapper.xml"})
	Init(option)

	if nil == conf {
		logger.LOG.Info("db config == nil")
		return
	}

	gb := Get("ds")

	var result *TUser
	err := gb.Select("userMapper.findById", map[string]interface{}{
		"id": 2,
	})(&result)

	fmt.Println("result:", result, "err:", err)

	var result2 *TUser
	err = gb.SelectContext(context.Background(), "userMapper.findById", map[string]interface{}{
		"id": 4,
	})(&result2)
	fmt.Println("result:", result2, "err:", err)
}

func TestGoBatisWithCodeConf(t *testing.T) {
	ds1 := NewDataSourceBuilder().
		DataSource("ds1").
		DriverName("mysql").
		DataSourceName("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8").
		MaxLifeTime(120).
		MaxOpenConns(10).
		MaxIdleConns(5).
		Build()

	option := NewDSOption().
		DS([]*DataSource{ds1}).
		Mappers([]string{"examples/mapper/userMapper.xml"}).
		ShowSQL(true)
	Init(option)

	if nil == conf {
		logger.LOG.Error("db config == nil")
		return
	}

	gb := Get("ds1")

	//result := make(map[string]interface{})
	//result := make([]interface{}, 0)
	//var result interface{}
	//result := make([]TUser, 0)
	var result *TUser
	err := gb.Select("userMapper.findById", map[string]interface{}{
		"id": 2,
	})(&result)

	fmt.Println("result:", result, "err:", err)

	u := &TUser{
		Name: "wenj1991",
		Password: na.NullString{
			String: "654321",
			Valid:  true,
		},
	}

	id, a, err := gb.Insert("userMapper.saveUser", u)
	fmt.Println("id:", id, "affected:", a, "err:", err)

	uu := &TUser{
		Id:   1,
		Name: "wenj1993",
		Password: na.NullString{
			String: "654321",
			Valid:  true,
		},
	}

	// test set
	affected, err := gb.Update("userMapper.updateByCond", uu)
	fmt.Println("updateByCond:", affected, err)

	param := &TUser{
		Name: "wenj1993",
	}

	// test where
	res := make([]*TUser, 0)
	err = gb.Select("userMapper.queryStructsByCond", param)(&res)
	fmt.Println("queryStructsByCond", res, err)

	// test trim
	res2 := make([]*TUser, 0)
	err = gb.Select("userMapper.queryStructsByCond2", param)(&res2)
	fmt.Println("queryStructsByCond", res2, err)

	affected, err = gb.Delete("userMapper.deleteById", map[string]interface{}{
		"id": 3,
	})
	fmt.Println("delete affected:", affected, "err:", err)
}
