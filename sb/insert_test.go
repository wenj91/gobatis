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

	query := Model(&c).
		Insert().
		Build()

	expectedQuery := "insert into customer (id, name, phone) values (#{ID}, #{Name}, #{Phone})"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}

func TestInsertWithNil(t *testing.T) {
	c := customer4{
		ID:    0,
		Name:  "Name",
		Phone: nil,
	}

	query := Model(&c).
		Insert().
		Build()

	expectedQuery := "insert into customer (id, name) values (#{ID}, #{Name})"
	assert.True(t, query == expectedQuery, fmt.Sprintf("bad query: %s", query))
}
