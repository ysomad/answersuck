package domain

import "time"

type Question struct {
	Id          int       `json:"id"`
	Q           string    `json:"question"`
	Answer      string    `json:"answer"`
	AnswerImage *string   `json:"answerImage"`
	Media       *string   `json:"media"`
	MediaType   *string   `json:"mediaType"`
	Author      string    `json:"author"`
	LanguageId  int       `json:"languageId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
