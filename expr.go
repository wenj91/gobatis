package gobatis

import (
	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/datasource"
	"github.com/araddon/qlbridge/vm"
  "github.com/araddon/qlbridge/value"
)

func exprProcess(expression string, mapper map[string]interface{}) (value.Value, bool)  {
	exprAst := expr.MustParse(expression)
	evalContext := datasource.NewContextSimpleNative(mapper)
	return vm.Eval(evalContext, exprAst)
}
