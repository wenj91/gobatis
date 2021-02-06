package sb

import (
	"github.com/wenj91/gobatis/m"
	"github.com/wenj91/gobatis/uti/field"
	"strings"
)

// Select returns a new SELECT statement with the default dialect.
func Select(cols ...string) SelectStatement {
	stmt := SelectStatement{}
	if len(cols) > 0 {
		stmt.selects = append(stmt.selects, cols...)
	}
	return stmt
}

// SelectStatement represents a SELECT statement.
type SelectStatement struct {
	model      m.Model
	resultType m.ResultType
	selects    []string
	joins      []join
	wheres     []Cond
	lock       bool
	limit      *int
	offset     *int
	orders     []Od
	group      []string
	having     string
}

type join struct {
	sql string
}

// Join returns a new statement with JOIN expression 'sql'.
func (s SelectStatement) RT(rt m.ResultType) SelectStatement {
	s.resultType = rt
	return s
}

// Join returns a new statement with JOIN expression 'sql'.
func (s SelectStatement) Join(sql string, sq ...string) SelectStatement {
	s.joins = append(s.joins, join{sql})
	for _, ss := range sq {
		s.joins = append(s.joins, join{ss})
	}
	return s
}

// Where returns a new statement with condition 'cond'. Multiple conditions
// are combined with AND.
func (s SelectStatement) Where(c Cond, cond ...Cond) SelectStatement {
	s.wheres = append(s.wheres, c)

	if len(cond) > 0 {
		for _, c := range cond {
			s.wheres = append(s.wheres, c)
		}
	}

	return s
}

// Limit returns a new statement with the limit set to 'limit'.
func (s SelectStatement) Limit(limit int) SelectStatement {
	s.limit = &limit
	return s
}

// Offset returns a new statement with the offset set to 'offset'.
func (s SelectStatement) Offset(offset int) SelectStatement {
	s.offset = &offset
	return s
}

// Od returns a new statement with ordering 'orders'.
// Only the last Od() is used.
func (s SelectStatement) Order(order Od, o ...Od) SelectStatement {
	s.orders = append(s.orders, order)
	if len(o) > 0 {
		s.orders = append(s.orders, o...)
	}

	return s
}

// Group returns a new statement with grouping 'group'.
// Only the last Group() is used.
func (s SelectStatement) Group(group string, gr ...string) SelectStatement {
	s.group = append(s.group, group)
	if len(gr) > 0 {
		s.group = append(s.group, gr...)
	}

	return s
}

// Having returns a new statement with HAVING condition 'having'.
// Only the last Having() is used.
func (s SelectStatement) Having(having string) SelectStatement {
	s.having = having
	return s
}

// Lock returns a new statement with FOR UPDATE locking.
func (s SelectStatement) Lock() SelectStatement {
	s.lock = true
	return s
}

// Build builds the SQL query. It returns the query, the argument slice,
// and the destination slice.
func (s SelectStatement) Build() (query string, args []interface{}) {
	if nil == s.model {
		panic("model must not be nil")
	}

	var cols []string
	if len(s.selects) > 0 {
		cols = append(cols, s.selects...)
	} else {
		fields := field.Fields(s.model)
		if len(fields) > 0 {
			cols = append(cols, fields...)
		} else {
			cols = append(cols, "*")
		}
	}
	query = "select " + strings.Join(cols, ", ") + " from " + s.model.Table()

	for _, join := range s.joins {
		sql := join.sql
		query += " " + sql
	}

	if len(s.wheres) > 0 {
		ss, vals := buildCond(s.wheres)
		query += ss
		args = append(args, vals...)
	}

	if len(s.orders) > 0 {
		ods := make([]string, 0)
		for _, o := range s.orders {
			ods = append(ods, o.expr())
		}

		if len(ods) > 0 {
			query += " order by " + strings.Join(ods, ", ")
		}
	}

	if len(s.group) > 0 {
		query += " group by " + strings.Join(s.group, ", ")
	}

	if s.having != "" {
		query += " having " + s.having
	}

	if s.limit != nil {
		query += " limit ?"
		args = append(args, *s.limit)
	}

	if s.offset != nil {
		query += " offset ?"
		args = append(args, *s.offset)
	}

	if s.lock {
		query += " for update"
	}

	return
}
