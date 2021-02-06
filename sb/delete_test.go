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
	query := Model(&c).
		Delete().
		Build()

	expectedQuery := "delete from customer"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestDeleteWhere(t *testing.T) {
	c := customer3{}
	query := Model(&c).
		Delete().
		Where(
			Eq("name", "Name"),
			Eq("phone", "Phone"),
		).
		Build()

	expectedQuery := "delete from customer where (name = #{Name}) and (phone = #{Phone})"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}
