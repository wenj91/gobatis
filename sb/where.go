package sb

import (
	"fmt"
	"strings"
)

func buildCond(cs []Cond) string {
	sqls := make([]string, 0)
	for _, where := range cs {
		sql := "(" + where.expr() + ")"
		sqls = append(sqls, sql)
	}
	if len(sqls) > 0 {
		return " where " + strings.Join(sqls, " and ")
	}

	return ""
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
	expr() string
}

type defaultCond struct {
	operator op
	field    string
	mark     string
	btStart  string
	btEnd    string
}

func (e defaultCond) expr() string {
	switch e.operator {
	case eq:
		return fmt.Sprintf("%s = #{%s}", e.field, e.mark)
	case ne:
		return fmt.Sprintf("%s != #{%s}", e.field, e.mark)
	case gt:
		return fmt.Sprintf("%s > #{%s}", e.field, e.mark)
	case ge:
		return fmt.Sprintf("%s >= #{%s}", e.field, e.mark)
	case lt:
		return fmt.Sprintf("%s < #{%s}", e.field, e.mark)
	case le:
		return fmt.Sprintf("%s <= #{%s}", e.field, e.mark)
	case between:
		return fmt.Sprintf("%s between #{%s} and #{%s}", e.field, e.btStart, e.btEnd)
	case notBetween:
		return fmt.Sprintf("%s not between #{%s} and #{%s}", e.field, e.btStart, e.btEnd)
	case like:
		return fmt.Sprintf("%s like concat('%%', #{%s}, '%%')", e.field, e.mark)
	case notLike:
		return fmt.Sprintf("%s not like concat('%%', #{%s}, '%%')", e.field, e.mark)
	case likeLeft:
		return fmt.Sprintf("%s not like concat('%%', #{%s})", e.field, e.mark)
	case likeRight:
		return fmt.Sprintf("%s not like concat(#{%s}, '%%')", e.field, e.mark)
	case isNull:
		return fmt.Sprintf("%s is null", e.field)
	case isNotNull:
		return fmt.Sprintf("%s is not null", e.field)
	case in:
		return fmt.Sprintf(`<foreach collection="%s" item="item" index="index" open="%s in (" close=")" separator=",">#{item}</foreach>`, e.mark, e.field)
	case notIn:
		return fmt.Sprintf(`<foreach collection="%s" item="item" index="index" open="%s not in (" close=")" separator=",">#{item}</foreach>`, e.mark, e.field)
	default:

	}
	return ""
}

func Eq(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: eq,
		field:    field,
		mark:     mark,
	}
}

func Ne(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: ne,
		field:    field,
		mark:     mark,
	}
}

func Gt(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: gt,
		field:    field,
		mark:     mark,
	}
}

func Ge(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: ge,
		field:    field,
		mark:     mark,
	}
}

func Lt(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: lt,
		field:    field,
		mark:     mark,
	}
}

func Le(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: le,
		field:    field,
		mark:     mark,
	}
}

func Between(field, btStart, btEnd string, cond ...bool) Cond {
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

func NotBetween(field, btStart, btEnd string, cond ...bool) Cond {
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

func Like(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: like,
		field:    field,
		mark:     mark,
	}
}

func NotLike(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: notLike,
		field:    field,
		mark:     mark,
	}
}

func LikeLeft(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: likeLeft,
		field:    field,
		mark:     mark,
	}
}

func LikeRight(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: likeRight,
		field:    field,
		mark:     mark,
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

func In(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: in,
		field:    field,
		mark:     mark,
	}
}

func NotIn(field, mark string, cond ...bool) Cond {
	if len(cond) > 0 {
		b := cond[0]
		if !b {
			return nil
		}
	}

	return &defaultCond{
		operator: notIn,
		field:    field,
		mark:     mark,
	}
}
