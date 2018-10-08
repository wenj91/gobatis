package gobatis

import (
	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/datasource"
	"github.com/araddon/qlbridge/vm"
)

func exprProcess(expression string, mapper map[string]interface{})  bool  {
	exprAst := expr.MustParse(expression)
	evalContext := datasource.NewContextSimpleNative(mapper)
	v, _ := vm.Eval(evalContext, exprAst)
	return v.Value().(bool)
}
