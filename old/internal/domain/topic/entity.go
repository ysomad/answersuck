package topic

import (
	"errors"
	"time"

	"github.com/ysomad/answersuck-backend/internal/pkg/pagination"
)

var (
	ErrLanguageNotFound = errors.New("language with given id not found")
)

type Topic struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	LanguageId uint      `json:"language_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Filter struct {
	Name       string
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
