package sort

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Sort struct {
	Column string
	Order  string
}

func New(column, order string) Sort {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	return Sort{
		Column: column,
		Order:  order,
	}
}

// UseSelectBuilder adds sort to squirrel.SelectBuilder
func (s Sort) UseSelectBuilder(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(s.Column + " " + s.Order)
}
