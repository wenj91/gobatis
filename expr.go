package gobatis

import (
	"github.com/antonmedv/expr"
	"github.com/wenj91/gobatis/logger"
)

func eval(expression string, mapper map[string]interface{}) bool {
	ok, err := expr.Eval(expression, mapper)
	if nil != err {
		logger.LOG.Debug("[WARN]", "Expression:", expression, ">>> eval result err:", err)
		return false
	}

	return ok.(bool)
}
