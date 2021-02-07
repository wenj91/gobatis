package gobatis

import (
	"context"
	"github.com/wenj91/gobatis/sb"
)

type Wrapper struct {
	rt string
	gb *gbBase
	wp sb.Wrapper
}

func (w Wrapper) ResultType(rt string) Wrapper {
	w.rt = rt
	return w
}

func (w Wrapper) Select(fn func(s *sb.SelectStatement)) func(res interface{}) error {
	if w.rt == "" {
		panic("query result type must not be nil")
	}

	statement := w.wp.Select()
	fn(statement)

	return func(res interface{}) error {
		executor := &executor{
			gb: w.gb,
		}

		sqlStr, paramArr := statement.Build()
		err := executor.wrapperQuery(sqlStr, ResultType(w.rt), nil, paramArr, res)
		return err
	}
}

func (w Wrapper) SelectContext(ctx context.Context, fn func(s *sb.SelectStatement)) func(res interface{}) error {
	if w.rt == "" {
		panic("query result type must not be nil")
	}

	statement := w.wp.Select()
	fn(statement)

	return func(res interface{}) error {
		executor := &executor{
			gb: w.gb,
		}

		sqlStr, paramArr := statement.Build()
		err := executor.wrapperQueryContext(ctx, sqlStr, ResultType(w.rt), nil, paramArr, res)
		return err
	}
}

func (w Wrapper) Insert() (lastInsertId int64, affected int64, err error) {
	executor := &executor{
		gb: w.gb,
	}

	sqlStr, paramArr := w.wp.Insert().Build()
	return executor.wrapperUpdate(sqlStr, nil, paramArr)
}

func (w Wrapper) InsertContext(ctx context.Context) (lastInsertId int64, affected int64, err error) {
	executor := &executor{
		gb: w.gb,
	}

	sqlStr, paramArr := w.wp.Insert().Build()
	return executor.wrapperUpdateContext(ctx, sqlStr, nil, paramArr)
}

func (w Wrapper) Update(fn func(s *sb.UpdateStatement)) (affected int64, err error) {
	executor := &executor{
		gb: w.gb,
	}

	statement := w.wp.Update()
	fn(statement)

	sqlStr, paramArr := statement.Build()
	_, affected, err = executor.wrapperUpdate(sqlStr, nil, paramArr)
	return
}

func (w Wrapper) UpdateContext(ctx context.Context, fn func(s *sb.UpdateStatement)) (affected int64, err error) {
	executor := &executor{
		gb: w.gb,
	}

	statement := w.wp.Update()
	fn(statement)

	sqlStr, paramArr := statement.Build()
	_, affected, err = executor.wrapperUpdateContext(ctx, sqlStr, nil, paramArr)
	return
}

func (w Wrapper) Delete(fn func(s *sb.DeleteStatement)) (affected int64, err error) {
	executor := &executor{
		gb: w.gb,
	}

	statement := w.wp.Delete()
	fn(statement)

	sqlStr, paramArr := statement.Build()
	_, affected, err = executor.wrapperUpdate(sqlStr, nil, paramArr)
	return
}

func (w Wrapper) DeleteContext(ctx context.Context, fn func(s *sb.DeleteStatement)) (affected int64, err error) {
	executor := &executor{
		gb: w.gb,
	}

	statement := w.wp.Delete()
	fn(statement)

	sqlStr, paramArr := statement.Build()
	_, affected, err = executor.wrapperUpdateContext(ctx, sqlStr, nil, paramArr)
	return
}
