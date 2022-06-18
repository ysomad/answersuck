package question

import "time"

type Question struct {
	Id         int       `json:"id"`
	Text       string    `json:"text"`
	AnswerId   int       `json:"answerId"`
	MediaId    *string   `json:"mediaId"`
	AccountId  string    `json:"accountId"`
	LanguageId int       `json:"languageId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (q *Question) PrepareForSave() {
	now := time.Now()
	q.CreatedAt = now
	q.UpdatedAt = now
}

// Detailed is question entity with joined tables associated with it
type Detailed struct {
	Id             int       `json:"id"`
	Text           string    `json:"text"`
	Answer         string    `json:"answer"`
	AnswerMediaURL *string   `json:"answerMediaUrl"`
	Author         string    `json:"author"`
	MediaURL       *string   `json:"media"`
	MediaType      *string   `json:"mediaType"`
	LanguageId     int       `json:"languageId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Minimized is minimized question entity using for lists
type Minimized struct {
	Id         int    `json:"id"`
	Text       string `json:"text"`
	LanguageId int    `json:"languageId"`
}
