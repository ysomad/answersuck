package question

import "time"

type Question struct {
	Id             int       `json:"id"`
	Text           string    `json:"text"`
	Answer         string    `json:"answer"`
	AnswerImageURL *string   `json:"answerImageUrl"`
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
