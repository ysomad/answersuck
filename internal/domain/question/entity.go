package question

import (
	"errors"
	"time"
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
	Id         uint32 `json:"id"`
	Text       string `json:"text"`
	LanguageId uint8  `json:"language_id"`
}
