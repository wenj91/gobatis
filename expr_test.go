package gobatis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		"arr":  []string{"1", "2"},
		"arr2": []string{},
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
		"val != nil and val == ''",
		"arr != nil and len(arr) > 0",
		"arr2 != nil and len(arr2) > 0",
		"$blank(val)",
		"!$blank(val)",
	}

	for i, ex := range expression {
		ok := eval(ex, params)
		fmt.Printf("Index:%v Expr:%v >>>> Result:%v \n", i, ex, ok)
		assertExpr(t, i, ok, ex)
	}
}

func assertExpr(t *testing.T, i int, ok bool, expr string) {
	switch i {
	case 0, 3, 4, 6, 8, 9, 12, 14, 17, 22, 24: // false
		assert.True(t, !ok, "Expr:"+expr+" Result:true")
	case 1, 2, 5, 7, 10, 11, 13, 15, 16, 18, 19, 20, 21, 23: // true
		assert.True(t, ok, "Expr:"+expr+" Result:false")
	}
}
