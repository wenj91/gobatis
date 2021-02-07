package sb

import (
	"fmt"
	"github.com/wenj91/gobatis/m"
	"strings"
)

type updateSet struct {
	col string
	arg interface{}
	raw bool
}

// *UpdateStatement represents an update statement.
type UpdateStatement struct {
	model  m.Model
	sets   []updateSet
	wheres []cond
}

// Set returns a new statement with column 'col' set to value 'val'.
func (s *UpdateStatement) Set(col string, arg interface{}) *UpdateStatement {
	s.sets = append(s.sets, updateSet{col: col, arg: arg, raw: false})
	return s
}

// Where returns a new statement with condition 'cond'.
// Multiple Where() are combined with and.
func (s *UpdateStatement) Where(cond cond, cs ...cond) *UpdateStatement {
	s.wheres = append(s.wheres, cond)
	if len(cs) > 0 {
		for _, c := range cs {
			s.wheres = append(s.wheres, c)
		}
	}
	return s
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s *UpdateStatement) Build() (query string, args []interface{}) {
	if len(s.sets) == 0 {
		panic("sqlbuilder: no columns set")
	}

	query = "update " + s.model.Table() + " set "
	var sets []string

	for _, set := range s.sets {
		sets = append(sets, fmt.Sprintf("%s = ?", set.col))
		args = append(args, set.arg)
	}
	query += strings.Join(sets, ", ")

	if len(s.wheres) > 0 {
		ss, v := buildCond(s.wheres)
		query += ss
		args = append(args, v...)
	}

	return
}
