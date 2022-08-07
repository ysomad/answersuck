package filter

import sq "github.com/Masterminds/squirrel"

type Type uint8

const (
	// TypeEQ value is equal
	TypeEQ Type = iota

	// TypeNotEQ value is not equal
	TypeNotEQ

	// TypeGTE greater than value
	TypeGTE

	// TypeGT value greater
	TypeGT

	// TypeLT value less
	TypeLT

	// TypeLTE value less or equals
	TypeLTE

	// TypeLike value can contain
	TypeLike

	// TypeNotLike value cannot contain
	TypeNotLike

	// TypeILike value can contain (case insensitive)
	TypeILike

	// TypeNotILike value cannot contain (case insensitive)
	TypeNotILike
)

type Operator uint8

const (
	OperatorAnd Operator = iota
	OperatorOr
)

type Filter struct {
	column   string
	ftype    Type
	value    any
	operator Operator
	filters  []Filter
}

func New(column string, ftype Type, value any) Filter {
	return Filter{
		column:   column,
		ftype:    ftype,
		value:    value,
		operator: OperatorAnd,
		filters:  make([]Filter, 0),
	}
}

// SetOperator sets operator for linking filters
func (f *Filter) SetOperator(operator Operator) *Filter {
	f.operator = operator
	return f
}

// WithFilter adds filters
func (f *Filter) WithFilters(filters ...Filter) *Filter {
	f.filters = append(f.filters, filters...)
	return f
}

func (f Filter) condition() sq.Sqlizer {
	switch f.ftype {
	case TypeEQ:
		return sq.Eq{f.column: f.value}
	case TypeNotEQ:
		return sq.NotEq{f.column: f.value}
	case TypeGTE:
		return sq.GtOrEq{f.column: f.value}
	case TypeGT:
		return sq.Gt{f.column: f.value}
	case TypeLT:
		return sq.Lt{f.column: f.value}
	case TypeLTE:
		return sq.LtOrEq{f.column: f.value}
	case TypeLike:
		return sq.Like{f.column: f.value}
	case TypeNotLike:
		return sq.NotLike{f.column: f.value}
	case TypeILike:
		return sq.ILike{f.column: f.value}
	case TypeNotILike:
		return sq.NotILike{f.column: f.value}
	}
	return sq.Eq{f.column: f.value}
}

func (f Filter) getConditions() sq.Sqlizer {
	if len(f.filters) == 0 {
		return f.condition()
	}

	var conditions []sq.Sqlizer

	conditions = append(conditions, f.condition())

	for _, filter := range f.filters {
		conditions = append(conditions, filter.getConditions())
	}

	if f.operator == OperatorOr {
		return or(conditions)
	}

	return and(conditions)
}

// UseSelectBuilder adds filters to squirrel.SelectBuilder
func (f Filter) UseSelectBuilder(b sq.SelectBuilder) sq.SelectBuilder {
	return b.Where(f.getConditions())
}

func and(conditions []sq.Sqlizer) sq.And {
	res := sq.And{}
	for _, c := range conditions {
		res = append(res, c)
	}
	return res
}

func or(conditions []sq.Sqlizer) sq.Or {
	res := sq.Or{}
	for _, c := range conditions {
		res = append(res, c)
	}
	return res
}
