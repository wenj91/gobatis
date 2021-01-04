package examples

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 引入驱动
	"github.com/wenj91/gobatis"        // 引入gobatis
)

// 实体结构示例， tag：field为数据库对应字段名称
type User struct {
	Id    gobatis.NullInt64  `field:"id"`
	Name  gobatis.NullString `field:"name"`
	Email gobatis.NullString `field:"email"`
	CrtTm gobatis.NullTime   `field:"crtTm"`
}

func main() {
	// 初始化db
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	dbs := make(map[string]*gobatis.GoBatisDB)
	dbs["ds1"] = gobatis.NewGoBatisDB(gobatis.DBTypeMySQL, db)

	option := gobatis.NewDBOption().
		DB(dbs).
		ShowSQL(true).
		Mappers([]string{"mapper/userMapper.xml"})

	gobatis.Init(option)

	// 获取数据源，参数为数据源名称，如：ds1
	gb := gobatis.Get("ds1")

	//传入id查询Map
	mapRes := make(map[string]interface{})
	// stmt标识为：namespace + '.' + id, 如：userMapper.findMapById
	// 查询参数可以是map，也可以是数组，也可以是实体结构
	err := gb.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(mapRes)
	fmt.Println("userMapper.findMapById-->", mapRes, err)

	mapRes2 := make(map[string]interface{})
	err = gb.SelectContext(context.TODO(), "userMapper.findMapById", map[string]interface{}{"id": 4})(mapRes2)
	fmt.Println("userMapper.findMapById-->", mapRes2, err)
}
