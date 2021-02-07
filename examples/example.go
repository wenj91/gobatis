package examples

import (
	"database/sql"
	"fmt"
	"github.com/wenj91/gobatis/na"
	"github.com/wenj91/gobatis/sb"
	"github.com/wenj91/gobatis/uti"
	"time"

	_ "github.com/go-sql-driver/mysql" // 引入驱动
	"github.com/wenj91/gobatis"        // 引入gobatis
)

// 实体结构示例， tag：field为数据库对应字段名称
type User2 struct {
	Id    na.NullInt64  `field:"id"`
	Name  na.NullString `field:"name"`
	Email na.NullString `field:"email"`
	CrtTm na.NullTime   `field:"crtTm"`
}

func (u *User2) Table() string {
	return "user"
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

	u := &User2{
		Name:  uti.NS("wenj1991"),
		Email: uti.NS("email@cc.com"),
	}

	insertId, affected, err := gb.Wrapper(u).Insert()
	fmt.Printf("insert insertId:%d affected:%d err:%v\n", insertId, affected, err)

	m := make([]map[string]interface{}, 0)
	err = gb.Wrapper(u).ResultType("maps").Select(func(s *sb.SelectStatement) {
		s.Where(sb.Eq("id", 10))
	})(&m)
	fmt.Printf("select result:%v err:%v\n", m, err)

	a, err := gb.Wrapper(u).Delete(func(s *sb.DeleteStatement) {
		s.Where(sb.Eq("id", 9))
	})
	fmt.Printf("delete result:%v err:%v\n", a, err)

	i, err := gb.Wrapper(u).Update(func(s *sb.UpdateStatement) {
		s.Where(sb.Eq("id", 11)).
			Set("crtTm", time.Now())
	})
	fmt.Printf("update result:%v err:%v\n", i, err)
}
