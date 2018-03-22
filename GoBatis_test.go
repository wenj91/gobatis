package gobatis

import (
	"testing"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wenj91/gobatis"
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
	Id int
	Name string
	Email string
	CrtTm structs.NullTime
}

func TestNewGoBatisGetMap(t *testing.T) {
	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	//传入id查询Map
	mapRes := make(map[string]interface{})
	i, err := gobatis.Select("Mapper.findMapById", 1)(mapRes)
	fmt.Println("Mapper.findMapById-->", i, mapRes, err)
}

func TestNewGoBatisGetStruct(t *testing.T) {
	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	//根据传入实体查询对象
	param := User{
		Id: structs.NullInt64{1, true},
	}
	structRes2 := User{}
	i, err := gobatis.Select("Mapper.findStructByStruct", param)(&structRes2)
	fmt.Println("Mapper.findStructByStruct-->", i, structRes2, err)
}

func TestNewGoBatisInsertBatch(t *testing.T) {

	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	param := Test{
		Name: "wenj91",
		Email: "123@xxx.com",
		CrtTm: structs.NullTime{time.Now(), true},
	}
	param2 := Test{
		Name: "wenj92",
		Email: "123@xxx.com",
		CrtTm: structs.NullTime{time.Now(), true},
	}

	pp := []Test{param, param2}
	i, j, _ :=gobatis.Insert("Mapper.insertStructsBatch", pp)
	fmt.Println(i, j)
}

func TestNewGoBatisUpdate(t *testing.T) {
	//updateByStruct
	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	param := Test{
		Id: 1,
		Name: "wenj91update",
		Email: "123update@xxx.com",
		CrtTm: structs.NullTime{time.Now(), true},
	}

	i, err := gobatis.Update("Mapper.updateByStruct", param)
	fmt.Println("updateByStruct", i, err)
}

func TestNewGoBatisDelete(t *testing.T) {
	//deleteById
	mapperPath := []string{"./mapper.xml"}
	gobatis := gobatis.NewGoBatis("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", mapperPath)

	i, err := gobatis.Delete("Mapper.deleteById", 1011)
	fmt.Println("deleteById", i, err)
}
