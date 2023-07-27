package paging

import (
	"errors"
)

var errInvalidArgs = errors.New("paging: length of items should not equal more than pageSize + 1")

type List[T any] struct {
	Items         []T
	NextPageToken string
}

func NewListWithOffset[T any](items []T, pageSize, offset uint64) (List[T], error) {
	itemsLen := uint64(len(items))

	// more items than len(items) + 1
	if itemsLen > pageSize && itemsLen-pageSize != 1 {
		return List[T]{}, errInvalidArgs
	}

	// has no next page
	if itemsLen != pageSize+1 {
		return List[T]{
			Items:         items,
			NextPageToken: "",
		}, nil
	}

	return List[T]{
		Items:         items[:itemsLen-1],
		NextPageToken: string(NewOffsetToken(pageSize, offset+pageSize)),
	}, nil
}
