package gobatis

import (
	"testing"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"../gobatis"
	"github.com/wenj91/gobatis/structs"
	"time"
)

type User struct {
	Id    structs.NullInt64
	Name  structs.NullString
	Email structs.NullString
	CrtTm structs.NullTime
}

type Test struct {
	Name string
	Email string
	CrtTm structs.NullTime
}

func TestNewGoBatis(t *testing.T) {
	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	//传入id查询Map
	//mapRes := make(map[string]interface{})
	//i, err := gobatis.Select("Mapper.findMapById", 1)(mapRes)
	//fmt.Println("Mapper.findMapById-->", i, mapRes, err)

	//根据传入实体查询对象
	param := User{
		Id: structs.NullInt64{1, true},
	}
	structRes2 := User{}
	i, err := gobatis.Select("Mapper.findStructByStruct", param)(&structRes2)
	fmt.Println("Mapper.findStructByStruct-->", i, structRes2, err)
}

func TestNewGoBatis2(t *testing.T) {

	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	param := Test{
		"wenj91",
		"123@xxx.com",
		structs.NullTime{time.Now(), true},
	}
	param2 := Test{
		"wenj92",
		"123@xxx.com",
		structs.NullTime{time.Now(), true},
	}
	param3 := Test{
		"wenj93",
		"123@xxx.com",
		structs.NullTime{time.Now(), true},
	}

	pp := []Test{param, param2, param3}
	i, j, _ :=gobatis.Insert("Mapper.insertStructsBatch", pp)
	fmt.Println(i, j)
}