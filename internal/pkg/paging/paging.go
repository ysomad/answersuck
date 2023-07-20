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

func NewList[T any](items []T, pageSize uint64) (List[T], error) {
	itemsLen := uint64(len(items))

	if itemsLen > pageSize && itemsLen-pageSize != 1 {
		return List[T]{}, errInvalidArgs
	}

	if itemsLen != pageSize+1 {
		return List[T]{
			Items:   items,
			HasNext: false,
		}, nil
	}

	return List[T]{
		Items:   items[:itemsLen-1],
		HasNext: true,
	}, nil
}

// SeekParams is set of params for seek(keyset) pagination using uuids
// or other non-sortable primary keys.
type SeekParams struct {
	PageToken Token
	PageSize  int32
}

type SeekListItem interface {
	GetID() string
	GetTime() time.Time
}

type SeekList[T SeekListItem] struct {
	Items         []T
	NextPageToken Token
}

func NewSeekList[T SeekListItem](items []T, pageSize int32) (SeekList[T], error) {
	if len(items) > int(pageSize) && len(items)-int(pageSize) != 1 {
		return SeekList[T]{}, errInvalidArgs
	}

	if len(items) != int(pageSize+1) {
		return SeekList[T]{Items: items}, nil
	}

	list := items[:len(items)-1]

	return SeekList[T]{
		Items:         list,
		NextPageToken: NewToken(list[len(list)-1].GetID(), list[len(list)-1].GetTime()),
	}, nil
}

type OffsetParams struct {
	Offset uint64
	Limit  uint64
}
