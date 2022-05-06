package domain

import "time"

type Question struct {
	Id             int       `json:"id"`
	Q              string    `json:"question"`
	MediaURL       *string   `json:"media"`
	MediaType      *string   `json:"mediaType"`
	Answer         string    `json:"answer"`
	AnswerImageURL *string   `json:"answerImage"`
	Author         string    `json:"author"`
	LanguageId     int       `json:"languageId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
