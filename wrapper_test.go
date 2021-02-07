package gobatis

import (
	"fmt"
	"github.com/wenj91/gobatis/logger"
	"github.com/wenj91/gobatis/na"
	"github.com/wenj91/gobatis/sb"
	"github.com/wenj91/gobatis/uti"
	"testing"
	"time"
)

type TUser2 struct {
	Id       int64         `field:"id"`
	Name     string        `field:"name"`
	Password na.NullString `field:"pwd"`
	Email    na.NullString `field:"email"`
	CrtTm    na.NullTime   `field:"crtTm"`
}

func (u *TUser2) Table() string {
	return "user"
}

func TestGoBatisWrapper(t *testing.T) {
	Init(NewFileOption())
	if nil == conf {
		logger.LOG.Error("db config == nil")
		return
	}

	gb := Get("ds")

	u := &TUser2{
		Name:     "wenj1991",
		Password: uti.NS("password"),
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
