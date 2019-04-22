package gobatis

import (
	"log"

	"github.com/antonmedv/expr"
)

func eval(expression string, mapper map[string]interface{}) bool {
	ok, err := expr.Eval(expression, mapper)
	if nil != err {
		// panic here is better ??
		log.Println("[WARN]", "Expression:", expression, ">>> eval result err:", err)
		return false
	}

	return ok.(bool)
}
