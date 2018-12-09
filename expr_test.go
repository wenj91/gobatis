package gobatis

import (
	"fmt"
	"testing"
)

type TestUser struct {
	Name string
}

func TestExpr_eval(t *testing.T) {
	params := map[string]interface{}{
		"name": "wenj91",
		"val":  "",
		"user": &TestUser{Name: "wenj91"},
		"m":    map[string]interface{}{"user": &TestUser{Name: "wenj91"}},
		"m1":   map[string]interface{}{"name": "wenj91"},
	}
	expression := []string{
		"1 != 1",
		"1 == 1",
		"name == 'wenj91'",
		"name != 'wenj91'",
		"user.Name1 == 'wenj91'",
		"user.Name == 'wenj91'",
		"user.Name != 'wenj91'",
		"user.Name != nil",
		"user.Name == nil",
		"m.user.Name != 'wenj91'",
		"m.user.Name == 'wenj91'",
		"m1.name == 'wenj91'",
		"m1.name != 'wenj91'",
		"m.user.Name == 'wenj91' && 1 == 1",
		"m.user.Name == 'wenj91' && 1 != 1",
		"m.user.Name == 'wenj91' || 1 != 1",
		"val != nil",
		"val != ''",
		"val == ''",
		"val != nil && val == ''",
	}

	for i, ex := range expression {
		ok := eval(ex, params)
		fmt.Printf("Index:%v Expr:%v >>>> Result:%v \n", i, ex, ok)
		assertExpr(i, ok, ex)
	}
}

func assertExpr(i int, ok bool, expr string) {
	switch i {
	case 0, 3, 4, 6, 8, 9, 12, 14, 17: // false
		assertNotTrue(ok, "Expr:"+expr+" Result:true")
	case 1, 2, 5, 7, 10, 11, 13, 15, 16, 18, 19: // true
		assertTrue(ok, "Expr:"+expr+" Result:false")
	}
}
