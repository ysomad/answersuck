package answer

import (
	"errors"

	"github.com/answersuck/host/internal/pkg/pagination"
)

var (
	ErrMediaTypeNotAllowed = errors.New("not allowed media type for answer")
	ErrLanguageNotFound    = errors.New("language with provided id not found")
)

type Answer struct {
	Id         int     `json:"id"`
	Text       string  `json:"text"`
	MediaId    *string `json:"media_id"`
	LanguageId uint    `json:"language_id"`
}

type Filter struct {
	Text       string
	LanguageId uint
}

type ListParams struct {
	Pagination pagination.Params
	Filter     Filter
}

func NewListParams(lastId uint32, limit uint64, f Filter) ListParams {
	if limit == 0 || limit > pagination.MaxLimit {
		limit = pagination.DefaultLimit
	}
	return ListParams{
		Pagination: pagination.Params{
			LastId: lastId,
			Limit:  limit,
		},
		Filter: f,
	}
}
