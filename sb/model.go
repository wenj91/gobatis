package sb

import "github.com/wenj91/gobatis/m"

type Wrapper struct {
	model m.Model
}

func Model(m m.Model) Wrapper {
	return Wrapper{
		model: m,
	}
}

// Select returns a new SELECT statement with the default dialect.
func (dm Wrapper) Select(cols ...string) *SelectStatement {
	stmt := &SelectStatement{
		model: dm.model,
	}
	if len(cols) > 0 {
		stmt.selects = append(stmt.selects, cols...)
	}
	return stmt
}

// Update returns a new Update statement with the default dialect.
func (dm Wrapper) Update() *UpdateStatement {
	stmt := &UpdateStatement{
		model: dm.model,
	}
	return stmt
}

// Delete returns a new Delete statement with the default dialect.
func (dm Wrapper) Delete() *DeleteStatement {
	stmt := &DeleteStatement{
		model: dm.model,
	}
	return stmt
}

// Insert returns a new INSERT statement with the default dialect.
func (dm Wrapper) Insert() *InsertStatement {
	return &InsertStatement{
		model: dm.model,
	}
}
