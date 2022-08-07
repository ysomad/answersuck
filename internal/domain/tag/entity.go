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
	LanguageId uint   `json:"language_id"`
}

const maxLimit = 100

type Filter struct {
	Name       string
	LanguageId uint
}

type ListParams struct {
	Pagination pagination.Params
	Filter     Filter
}

func NewListParams(lastId uint32, limit uint64, f Filter) ListParams {
	if limit == 0 || limit > maxLimit {
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
