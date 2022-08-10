package question

import (
	"errors"
	"time"

	"github.com/answersuck/host/internal/pkg/pagination"
)

var (
	ErrForeignKeyViolation = errors.New("provided answer, media or language does not exist")
	ErrNotFound            = errors.New("question with provided id not found")
)

type Question struct {
	Id         uint32    `json:"id"`
	Text       string    `json:"text"`
	AnswerId   uint32    `json:"answer_id"`
	MediaId    *string   `json:"media_id"`
	AccountId  string    `json:"account_id"`
	LanguageId uint8     `json:"language_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Detailed is question entity with joined tables associated with it
type Detailed struct {
	Id              uint32    `json:"id"`
	Text            string    `json:"text"`
	Answer          string    `json:"answer"`
	AnswerMediaURL  *string   `json:"answer_media_url"`
	AnswerMediaType *string   `json:"answer_media_type"`
	Author          string    `json:"author"`
	MediaURL        *string   `json:"media_url"`
	MediaType       *string   `json:"media_type"`
	LanguageId      uint8     `json:"language_id"`
	CreatedAt       time.Time `json:"created_at"`
}

// setURLsFromFilenames in db storing media filenames and URL must be set
// manually before returning to end user
func (d *Detailed) setURLsFromFilenames(p mediaProvider) {
	switch {
	case d.AnswerMediaURL != nil:
		answerMediaURL := p.URL(*d.AnswerMediaURL).String()
		d.AnswerMediaURL = &answerMediaURL
	case d.MediaURL != nil:
		questionMediaURL := p.URL(*d.MediaURL).String()
		d.MediaURL = &questionMediaURL
	}
}

// Minimized is minimized question entity using for lists
type Minimized struct {
	Id         uint32    `json:"id"`
	Text       string    `json:"text"`
	Author     string    `json:"author"`
	LanguageId uint8     `json:"language_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Filter struct {
	Text       string
	Author     string
	LanguageId uint
}

type ListParams struct {
	Pagination pagination.Params
	Filter     Filter
	SortOrder  string
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
