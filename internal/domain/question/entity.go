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
	Id             uint32    `json:"id"`
	Text           string    `json:"text"`
	Answer         string    `json:"answer"`
	AnswerMediaURL *string   `json:"answer_media_url"`
	Author         string    `json:"author"`
	MediaURL       *string   `json:"media_url"`
	MediaType      *string   `json:"media_type"`
	LanguageId     uint8     `json:"language_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Minimized is minimized question entity using for lists
type Minimized struct {
	Id         uint32 `json:"id"`
	Text       string `json:"text"`
	LanguageId uint8  `json:"languageId"`
}
