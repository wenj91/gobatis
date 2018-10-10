package gobatis

import (
	"github.com/antonmedv/expr"
	"log"
)

func eval(expression string, mapper map[string]interface{}) bool {
	ok, err := expr.Eval(expression, mapper)
	if nil != err {
		// panic here or better ??
		log.Println(expression, "eval result err:", err)
		return false
	}

	return ok.(bool)
}
