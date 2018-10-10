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
	}

	for _, ex := range expression {
		ok := eval(ex, params)
		fmt.Printf("Eexpr:%v >>>> Result:%v \n", ex, ok)
	}
}
