package sb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customer3 struct {
	ID    int     `field:"id"`
	Name  string  `field:"name"`
	Phone *string `field:"phone"`
}

func (c *customer3) Table() string {
	return "customer"
}

func TestDelete(t *testing.T) {
	c := customer3{}
	query, args := Model(&c).
		Delete().
		Build()

	expectedQuery := "delete from customer"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 0)
}

func TestDeleteWhere(t *testing.T) {
	c := customer3{}
	query, args := Model(&c).
		Delete().
		Where(
			Eq("name", "Name"),
			Eq("phone", "Phone"),
		).
		Build()

	expectedQuery := "delete from customer where (name = ?) and (phone = ?)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 2)
	assert.True(t, args[0] == "Name", fmt.Sprintf("args[0] err, excpect:%s actual:%s", "Name", args[0]))
	assert.True(t, args[1] == "Phone", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Phone", args[1]))
}
