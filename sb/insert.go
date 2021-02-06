package sb

import (
	"github.com/wenj91/gobatis/m"
	"github.com/wenj91/gobatis/uti/field"
	"strings"
)

type insertSet struct {
	col string
	arg interface{}
	raw bool
}

type insertRet struct {
	sql  string
	dest interface{}
}

// InsertStatement represents an INSERT statement.
type InsertStatement struct {
	model m.Model
	rets  []insertRet
}

// Return returns a new statement with a RETURNING clause.
func (s InsertStatement) Return(col string, dest interface{}) InsertStatement {
	s.rets = append(s.rets, insertRet{sql: col, dest: dest})
	return s
}

// Build builds the SQL query. It returns the SQL query and the argument slice.
func (s InsertStatement) Build() (query string, args []interface{}) {
	q := make([]string, 0)
	cols, vals := field.Map(s.model)
	args = append(args, vals...)
	for _, _ = range cols {
		q = append(q, "?")
	}

	query = "insert into " + s.model.Table() + " (" + strings.Join(cols, ", ") + ") values (" + strings.Join(q, ", ") + ")"

	return
}
