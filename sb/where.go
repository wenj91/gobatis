package sb

import (
	"fmt"
	"strings"
)

func buildCond(cs []Cond) (string, []interface{}) {
	vals := make([]interface{}, 0)
	sqls := make([]string, 0)
	for _, where := range cs {
		s, v := where.expr()
		sql := "(" + s + ")"
		sqls = append(sqls, sql)
		vals = append(vals, v...)
	}
	if len(sqls) > 0 {
		return " where " + strings.Join(sqls, " and "), vals
	}

	return "", vals
}

type op int

const (
	eq op = iota
	ne
	gt
	ge
	lt
	le
	between
	notBetween
	like
	notLike
	likeLeft
	likeRight
	isNull
	isNotNull
	in
	notIn
)

// eq ne gt ge lt le between notBetween like notLike likeLeft likeRight isNull isNotNull in notIn
type Cond interface {
	expr() (expr string, params []interface{})
}

type defaultCond struct {
	operator op
	field    string
	val      interface{}
	btStart  string
	btEnd    string
	in       []interface{}
}

func (e defaultCond) expr() (expr string, params []interface{}) {
	switch e.operator {
	case eq:
		return fmt.Sprintf("%s = ?", e.field), []interface{}{e.val}
	case ne:
		return fmt.Sprintf("%s != ?", e.field), []interface{}{e.val}
	case gt:
		return fmt.Sprintf("%s > ?", e.field), []interface{}{e.val}
	case ge:
		return fmt.Sprintf("%s >= ?", e.field), []interface{}{e.val}
	case lt:
		return fmt.Sprintf("%s < ?", e.field), []interface{}{e.val}
	case le:
		return fmt.Sprintf("%s <= ?", e.field), []interface{}{e.val}
	case between:
		return fmt.Sprintf("%s between ? and ?", e.field), []interface{}{e.btStart, e.btEnd}
	case notBetween:
		return fmt.Sprintf("%s not between ? and ?", e.field), []interface{}{e.btStart, e.btEnd}
	case like:
		return fmt.Sprintf("%s like concat('%%', ?, '%%')", e.field), []interface{}{e.val}
	case notLike:
		return fmt.Sprintf("%s not like concat('%%', ?, '%%')", e.field), []interface{}{e.val}
	case likeLeft:
		return fmt.Sprintf("%s not like concat('%%', ?)", e.field), []interface{}{e.val}
	case likeRight:
		return fmt.Sprintf("%s not like concat(?, '%%')", e.field), []interface{}{e.val}
	case isNull:
		return fmt.Sprintf("%s is null", e.field), []interface{}{e.val}
	case isNotNull:
		return fmt.Sprintf("%s is not null", e.field), []interface{}{e.val}
	case in:
		qs := make([]string, 0)
		ps := make([]interface{}, 0)
		for _, v := range e.in {
			qs = append(qs, "?")
			ps = append(ps, v)
		}
		return fmt.Sprintf(`%s in (%s)`, e.field, strings.Join(qs, ",")), ps
	case notIn:
		qs := make([]string, 0)
		ps := make([]interface{}, 0)
		for _, v := range e.in {
			qs = append(qs, "?")
			ps = append(ps, v)
		}
		return fmt.Sprintf(`%s not in (%s)`, e.field, strings.Join(qs, ",")), ps
	default:

	}
	return "", []interface{}{}
}

func Eq(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: eq,
		field:    field,
		val:      val,
	}
}

func Ne(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: ne,
		field:    field,
		val:      val,
	}
}

func Gt(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: gt,
		field:    field,
		val:      val,
	}
}

func Ge(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: ge,
		field:    field,
		val:      val,
	}
}

func Lt(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: lt,
		field:    field,
		val:      val,
	}
}

func Le(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: le,
		field:    field,
		val:      val,
	}
}

func Between(field string, btStart, btEnd string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: between,
		field:    field,
		btStart:  btStart,
		btEnd:    btEnd,
	}
}

func NotBetween(field string, btStart, btEnd string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: notBetween,
		field:    field,
		btStart:  btStart,
		btEnd:    btEnd,
	}
}

func Like(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: like,
		field:    field,
		val:      val,
	}
}

func NotLike(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: notLike,
		field:    field,
		val:      val,
	}
}

func LikeLeft(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: likeLeft,
		field:    field,
		val:      val,
	}
}

func LikeRight(field string, val interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: likeRight,
		field:    field,
		val:      val,
	}
}

func IsNull(field string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: isNull,
		field:    field,
	}
}

func IsNotNull(field string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: isNotNull,
		field:    field,
	}
}

func In(field string, vals []interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: in,
		field:    field,
		in:       vals,
	}
}

func NotIn(field string, vals []interface{}, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: notIn,
		field:    field,
		in:       vals,
	}
}
