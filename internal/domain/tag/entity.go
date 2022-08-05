package tag

import (
	"errors"

	"github.com/answersuck/host/internal/pkg/pagination"
)

var (
	ErrLanguageIdNotFound = errors.New("language id not found")
	ErrEmptyTagList       = errors.New("empty tag list")
)

type Tag struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	LanguageId uint8  `json:"language_id"`
}

type ListParams struct {
	pagination.CursorParams
}

func NewListParams(lastId uint32, limit uint64) ListParams {
	if limit == 0 || limit > pagination.MaxLimit {
		limit = pagination.DefaultLimit
	}
	return ListParams{pagination.CursorParams{
		LastId: lastId,
		Limit:  limit,
	}}
}
