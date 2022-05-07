package gobatis

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

func blank(arg interface{}) bool {
	if nil == arg {
		return true
	}

	res := fmt.Sprint(arg)
	if res == "" {
		return true
	}

	if strings.TrimSpace(res) == "" {
		return true
	}

	return false
}

func eval(expression string, mapper map[string]interface{}) bool {
	env := map[string]interface{}{
		"$blank": blank,
	}

	for k, v := range mapper {
		env[k] = v
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		LOG.Debug("[WARN]", "Expression:", expression, ">>> Compile result err:", err)
		return false
	}

	ok, err := expr.Run(program, env)
	if err != nil {
		LOG.Debug("[WARN]", "Expression:", expression, ">>> eval result err:", err)
		return false
	}

	return ok.(bool)
}
