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

	query, args := Model(&c).
		Select().
		Build()
	expectedQuery := "select id, name, phone from customer"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 0)
}

func TestSimpleSelectWithWhere(t *testing.T) {
	c := customer{}

	query, args := Model(&c).
		Select().
		Where(
			Eq("name", "Name"),
			Ne("phone", "Phone"),
		).
		Build()
	expectedQuery := "select id, name, phone from customer where (name = ?) and (phone != ?)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 2, fmt.Sprintf("the len of args err, excpect:%d actual:%d", 2, len(args)))
	assert.True(t, args[0] == "Name", fmt.Sprintf("args[0] err, excpect:%s actual:%s", "Name", args[0]))
	assert.True(t, args[1] == "Phone", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Phone", args[1]))
}

func TestSimpleSelectWithWhereIn(t *testing.T) {
	c := customer{}

	query, args := Model(&c).
		Select().
		Where(
			Eq("name", "Name"),
			Ne("phone", "Phone"),
			In("phone", []interface{}{"Phones", "Phones2"}),
		).
		Build()
	expectedQuery := "select id, name, phone from customer where (name = ?) and (phone != ?) and (phone in (?,?))"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 4, fmt.Sprintf("the len of args err, excpect:%d actual:%d", 4, len(args)))
	assert.True(t, args[0] == "Name", fmt.Sprintf("args[0] err, excpect:%s actual:%s", "Name", args[0]))
	assert.True(t, args[1] == "Phone", fmt.Sprintf("args[1] err, excpect:%s actual:%s", "Phone", args[1]))
	assert.True(t, args[2] == "Phones", fmt.Sprintf("args[2] err, excpect:%s actual:%s", "Phones", args[2]))
	assert.True(t, args[3] == "Phones2", fmt.Sprintf("args[3] err, excpect:%s actual:%s", "Phones2", args[3]))
}

func TestSimpleSelectWithLimitOffset(t *testing.T) {
	c := customer{}

	query, args := Model(&c).
		Select().
		Limit(5).
		Offset(10).
		Build()

	expectedQuery := "select id, name, phone from customer limit ? offset ?"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 2)
	assert.True(t, args[0] == 5, fmt.Sprintf("args[0] err, excpect:%d actual:%d", 5, args[0]))
	assert.True(t, args[1] == 10, fmt.Sprintf("args[1] err, excpect:%d actual:%d", 10, args[1]))
}

func TestSimpleSelectWithJoins(t *testing.T) {
	c := customer{}

	query, args := Model(&c).
		Select().
		Join(
			"inner join orders on orders.customer_id = customers.id",
			"left join items on items.order_id = orders.id",
		).
		Build()

	expectedQuery := "select id, name, phone from customer inner join orders on orders.customer_id = customers.id left join items on items.order_id = orders.id"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 0)
}

func TestSelectWithWhere(t *testing.T) {
	c := customer{}

	query, args := Model(&c).
		Select().
		Where(
			Eq("id", "id"),
			IsNotNull("name"),
		).
		Build()

	expectedQuery := "select id, name, phone from customer where (id = #{id}) and (name is not null)"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 0)
}

func TestSelectWithGroup(t *testing.T) {
	c := customer{}
	query, args := Model(&c).
		Select("COUNT(*)").
		Group("city", "name").
		Build()
	expectedQuery := "select COUNT(*) from customer group by city, name"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 0)
}

func TestSelectWithOrder(t *testing.T) {
	c := customer{}
	query, args := Model(&c).
		Select("COUNT(*)").
		Group("city", "name").
		Order(
			OrderDesc("name"),
			OrderAsc("city"),
		).
		Build()
	expectedQuery := "select COUNT(*) from customer order by name desc, city asc group by city, name"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
	assert.True(t, len(args) == 0)
}
