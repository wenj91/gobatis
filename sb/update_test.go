package sb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customer2 struct {
	ID    int     `field:"id"`
	Name  string  `field:"name"`
	Phone *string `field:"phone"`
}

func (c *customer2) Table() string {
	return "customer"
}

func TestUpdate(t *testing.T) {
	c := customer2{}
	query, args := Model(&c).
		Update().
		Set("name", "Name").
		Set("phone", "Phone").
		Build()

	expectedQuery := "update customer set name = ?, phone = ?"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 2)
	assert.True(t, args[0] == "Name", fmt.Sprintf("args[0] err, excpect:%s actual:%s", "Name", args[0]))
	assert.True(t, args[1] == "Phone", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Phone", args[1]))
}

func TestUpdateWithWhere(t *testing.T) {
	c := customer2{}
	query, args := Model(&c).Update().
		Set("name", "Name").
		Set("phone", "Phone").
		Where(
			Eq("id", "Id"),
		).
		Build()

	expectedQuery := "update customer set name = ?, phone = ? where (id = ?)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 3)
	assert.True(t, args[0] == "Name", fmt.Sprintf("args[0] err, excpect:%s actual:%s", "Name", args[0]))
	assert.True(t, args[1] == "Phone", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Phone", args[1]))
	assert.True(t, args[2] == "Id", fmt.Sprintf("args[2] err, excpect:%s actual:%s", "Id", args[2]))
}
