package sb

import "fmt"

type od int

const (
	asc od = iota
	desc
)

type Od interface {
	expr() string
}

type defaultOrder struct {
	field string
	o     od
}

func (d defaultOrder) expr() string {
	if d.o == desc {
		return fmt.Sprintf(`%s %s`, d.field, "desc")
	}

	return fmt.Sprintf(`%s %s`, d.field, "asc")
}

func OrderAsc(field string) Od {
	return &defaultOrder{
		field: field,
		o:     asc,
	}
}

func OrderDesc(field string) Od {
	return &defaultOrder{
		field: field,
		o:     desc,
	}
}
