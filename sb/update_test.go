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
	c := customer{}
	query := Model(&c).
		Update().
		Set("name", "Name").
		Set("phone", "Phone").
		Build()

	expectedQuery := "update customer set name = #{Name}, phone = #{Phone}"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestUpdateWithWhere(t *testing.T) {
	c := customer{}
	query := Model(&c).Update().
		Set("name", "Name").
		Set("phone", "Phone").
		Where(
			Eq("id", "Id"),
		).
		Build()

	expectedQuery := "update customer set name = #{Name}, phone = #{Phone} where (id = #{Id})"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}
