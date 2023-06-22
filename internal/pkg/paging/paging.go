package paging

import (
	"errors"
	"time"

	"golang.org/x/exp/constraints"
)

const (
	minPageSize     = 1
	maxPageSize     = 500
	defaultPageSize = 100
)

// IntSeek is set of params for seek(keyset) pagination using integer as primary keys.
type IntSeek[T constraints.Integer] struct {
	LastID   T // id of last item in list of items
	PageSize int32
}

func NewIntSeek[T constraints.Integer](lastID T, psize int32) IntSeek[T] {
	s := IntSeek[T]{
		LastID: lastID,
	}

	if psize < minPageSize {
		s.PageSize = defaultPageSize

		return s
	}

	if psize > maxPageSize {
		s.PageSize = maxPageSize

		return s
	}

	s.PageSize = psize

	return s
}

// List is a list of objects for infinity scroll pagination.
type List[T any] struct {
	Items   []T
	HasNext bool
}

var (
	errInvalidArgs = errors.New("paging: length of items should not equal more than pageSize + 1")
)

func NewList[T any](items []T, pageSize int32) (List[T], error) {
	if len(items) > int(pageSize) && len(items)-int(pageSize) != 1 {
		return List[T]{}, errInvalidArgs
	}

	if len(items) != int(pageSize+1) {
		return List[T]{
			Items:   items,
			HasNext: false,
		}, nil
	}

	return List[T]{
		Items:   items[:len(items)-1],
		HasNext: true,
	}, nil
}

// Seek is set of params for seek(keyset) pagination using uuids
// or other non-sortable primary keys.
type Seek struct {
	PageToken Token
	PageSize  int32
}

func NewSeek(t Token, psize int32) Seek {
	s := Seek{PageToken: t}

	if psize < minPageSize {
		s.PageSize = defaultPageSize

		return s
	}

	if psize > maxPageSize {
		s.PageSize = maxPageSize

		return s
	}

	s.PageSize = psize

	return s
}

type ListItem interface {
	GetID() string
	GetTime() time.Time
}

type TokenList[T ListItem] struct {
	Items         []T
	NextPageToken Token
}

func NewTokenList[T ListItem](items []T, pageSize int32) (TokenList[T], error) {
	if len(items) > int(pageSize) && len(items)-int(pageSize) != 1 {
		return TokenList[T]{}, errInvalidArgs
	}

	if len(items) != int(pageSize+1) {
		return TokenList[T]{Items: items}, nil
	}

	list := items[:len(items)-1]

	return TokenList[T]{
		Items:         list,
		NextPageToken: NewToken(list[len(list)-1].GetID(), list[len(list)-1].GetTime()),
	}, nil
}
