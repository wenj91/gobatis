package sb

import (
	"github.com/wenj91/gobatis/m"
)

// DeleteStatement represents a DELETE statement.
type DeleteStatement struct {
	model  m.Model
	wheres []Cond
}

// Where returns a new statement with condition 'cond'.
// Multiple Where() are combined with AND.
func (s *DeleteStatement) Where(cond Cond, cs ...Cond) *DeleteStatement {
	s.wheres = append(s.wheres, cond)
	if len(cs) > 0 {
		for _, c := range cs {
			s.wheres = append(s.wheres, c)
		}
	}
	return s
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s *DeleteStatement) Build() (query string, args []interface{}) {
	query = "delete from " + s.model.Table()

	if len(s.wheres) > 0 {
		ss, v := buildCond(s.wheres)
		query += ss
		args = append(args, v...)
	}

	return
}
