package sb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wenj91/gobatis/uti"
	"testing"
)

type customer4 struct {
	ID    int     `field:"id"`
	Name  string  `field:"name"`
	Phone *string `field:"phone"`
}

func (c *customer4) Table() string {
	return "customer"
}

func TestInsert(t *testing.T) {
	c := customer4{
		ID:    0,
		Name:  "Name",
		Phone: uti.PS("hodss"),
	}

	query, args := Model(&c).
		Insert().
		Build()

	expectedQuery := "insert into customer (id, name, phone) values (?, ?, ?)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 3, fmt.Sprintf("the len of args err, excpect:%d actual:%d", 3, len(args)))
	assert.True(t, args[0] == 0, fmt.Sprintf("args[0] err, excpect:%d actual:%d", 0, args[0]))
	assert.True(t, args[1] == "Name", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Name", args[1]))
	assert.True(t, args[2] == "hodss", fmt.Sprintf("args[2] err, excpect:%s actual:%s", "hodss", args[2]))
}

func TestInsertWithNil(t *testing.T) {
	c := customer4{
		ID:    0,
		Name:  "Name",
		Phone: nil,
	}

	query, args := Model(&c).
		Insert().
		Build()

	expectedQuery := "insert into customer (id, name) values (?, ?)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 2, fmt.Sprintf("the len of args err, excpect:%d actual:%d", 2, len(args)))
	assert.True(t, args[0] == 0, fmt.Sprintf("args[0] err, excpect:%d actual:%d", 0, args[0]))
	assert.True(t, args[1] == "Name", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Name", args[1]))
}
