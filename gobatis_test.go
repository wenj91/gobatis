package gobatis

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
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

	gb := NewGobatis()

	res, _ := gb.db.Query("select 1")
	cols, _ := res.Columns()
	fmt.Println(cols)

	//result := make(map[string]interface{})
	//result := make([]interface{}, 0)
	//var result interface{}
	//result := make([]TUser, 0)
	var result TUser
	err := gb.Select("userMapper.findById", map[string]interface{}{
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

	id, err := gb.Insert("userMapper.saveUser", u)
	fmt.Println("id:", id, "err:", err)

	uu := &TUser{
		Id:   3,
		Name: "wenj1993",
		Password: NullString{
			String: "654321",
			Valid:  true,
		},
	}
	affected, err := gb.Update("userMapper.updateByStruct", uu)
	fmt.Println("affected:", affected, "err:", err)

	affected, err = gb.Delete("userMapper.deleteById", map[string]interface{}{
		"id": 3,
	})
	fmt.Println("delete affected:", affected, "err:", err)
}
