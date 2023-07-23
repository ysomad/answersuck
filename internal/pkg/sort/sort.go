package sort

import (
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const (
	OrderASC  = "ASC"
	OrderDESC = "DESC"
)

type Sort struct {
	col   string
	order string
}

func newSort(col, order string) Sort {
	order = strings.ToUpper(order)
	if order != OrderASC && order != OrderDESC {
		order = OrderASC
	}

	return Sort{
		col:   col,
		order: order,
	}
}

// Attach attaches sort to builder.
func (s Sort) Attach(b sq.SelectBuilder) sq.SelectBuilder {
	return b.OrderBy(s.col + " " + s.order)
}

var errInvalidSortString = errors.New("invalid sort string")

type SortList []Sort

// NewSortList parses sort string into slice of sorts.
// Sort string must follow pattern: '{column} {order},{column} {order}...', order is optional
// and ASC by default.
func NewSortList(s string) ([]Sort, error) {
	if len(s) == 0 {
		return nil, nil
	}

	orderbys := strings.Split(s, ",")
	sorts := make([]Sort, len(orderbys))

	for i, r := range orderbys {
		orderby := string(r)
		if strings.HasPrefix(orderby, " ") || strings.HasSuffix(orderby, " ") {
			return nil, errInvalidSortString
		}

		parts := strings.Split(orderby, " ")

		switch len(parts) {
		case 1:
			sorts[i] = Sort{
				col:   parts[0],
				order: OrderASC,
			}
		case 2:
			sorts[i] = newSort(parts[0], parts[1])
		default:
			return nil, errInvalidSortString
		}
	}

	return sorts, nil
}
