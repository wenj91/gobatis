package sb

import (
	"fmt"
	"github.com/wenj91/gobatis/m"
	"strings"
)

type updateSet struct {
	col  string
	mark string
	raw  bool
}

// UpdateStatement represents an update statement.
type UpdateStatement struct {
	model  m.Model
	sets   []updateSet
	wheres []Cond
}

// Set returns a new statement with column 'col' set to value 'val'.
func (s UpdateStatement) Set(col string, mark string) UpdateStatement {
	s.sets = append(s.sets, updateSet{col: col, mark: mark, raw: false})
	return s
}

// Where returns a new statement with condition 'cond'.
// Multiple Where() are combined with and.
func (s UpdateStatement) Where(cond Cond, cs ...Cond) UpdateStatement {
	s.wheres = append(s.wheres, cond)
	if len(cs) > 0 {
		for _, c := range cs {
			s.wheres = append(s.wheres, c)
		}
	}
	return s
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s UpdateStatement) Build() (query string) {
	if len(s.sets) == 0 {
		panic("sqlbuilder: no columns set")
	}

	query = "update " + s.model.Table() + " set "
	var sets []string

	for _, set := range s.sets {
		sets = append(sets, fmt.Sprintf("%s = #{%s}", set.col, set.mark))
	}
	query += strings.Join(sets, ", ")

	if len(s.wheres) > 0 {
		query += buildCond(s.wheres)
	}

	return
}
