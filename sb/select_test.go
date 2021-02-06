package sb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customer struct {
	ID    int     `field:"id"`
	Name  string  `field:"name"`
	Phone *string `field:"phone"`
}

func (c *customer) Table() string {
	return "customer"
}

func TestSimpleSelect(t *testing.T) {
	c := customer{}

	query := Model(&c).
		Select().
		Build()
	expectedQuery := "select id, name, phone from customer"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSimpleSelectWithWhere(t *testing.T) {
	c := customer{}

	query := Model(&c).
		Select().
		Where(
			Eq("name", "Name"),
			Ne("phone", "Phone"),
		).
		Build()
	expectedQuery := "select id, name, phone from customer where (name = #{Name}) and (phone != #{Phone})"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSimpleSelectWithWhereIn(t *testing.T) {
	c := customer{}

	query := Model(&c).
		Select().
		Where(
			Eq("name", "Name"),
			Ne("phone", "Phone"),
			In("phone", "Phones"),
		).
		Build()
	expectedQuery := "select id, name, phone from customer where (name = #{Name}) and (phone != #{Phone}) and (<foreach collection=\"Phones\" item=\"item\" index=\"index\" open=\"phone in (\" close=\")\" separator=\",\">#{item}</foreach>)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSimpleSelectWithLimitOffset(t *testing.T) {
	c := customer{}

	query := Model(&c).
		Select().
		Limit(5).
		Offset(10).
		Build()

	expectedQuery := "select id, name, phone from customer limit 5 offset 10"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSimpleSelectWithJoins(t *testing.T) {
	c := customer{}

	query := Model(&c).
		Select().
		Join(
			"inner join orders on orders.customer_id = customers.id",
			"left join items on items.order_id = orders.id",
		).
		Build()

	expectedQuery := "select id, name, phone from customer inner join orders on orders.customer_id = customers.id left join items on items.order_id = orders.id"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSelectWithWhere(t *testing.T) {
	c := customer{}

	query := Model(&c).
		Select().
		Where(
			Eq("id", "id"),
			IsNotNull("name"),
		).
		Build()

	expectedQuery := "select id, name, phone from customer where (id = #{id}) and (name is not null)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSelectWithGroup(t *testing.T) {
	c := customer{}
	query := Model(&c).
		Select("COUNT(*)").
		Group("city", "name").
		Build()
	expectedQuery := "select COUNT(*) from customer group by city, name"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestSelectWithOrder(t *testing.T) {
	c := customer{}
	query := Model(&c).
		Select("COUNT(*)").
		Group("city", "name").
		Order(
			OrderDesc("name"),
			OrderAsc("city"),
		).
		Build()
	expectedQuery := "select COUNT(*) from customer order by name desc, city asc group by city, name"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}
