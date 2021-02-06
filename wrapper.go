package gobatis

import (
	"github.com/wenj91/gobatis/m"
	"github.com/wenj91/gobatis/sb"
)

type Wrapper struct {
	rt m.ResultType
	gb *gbBase
	wp sb.Wrapper
}

func (w Wrapper) Select(s sb.SelectStatement) func(res interface{}) error {
	return func(res interface{}) error {
		executor := &executor{
			gb: w.gb,
		}

		sqlStr, paramArr := s.Build()

		err := executor.wrapperQuery(sqlStr, paramArr, res)
		return err
	}
}

func (w Wrapper) Select(s sb.SelectStatement) func(res interface{}) error {
	return func(res interface{}) error {
		executor := &executor{
			gb: w.gb,
		}
		err := executor.w(ms, params, res)
		return err
	}
}
