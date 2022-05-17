package question

import "time"

type (
	CreateRequest struct {
		Question      string `json:"question" binding:"required,gte=4,lte=200"`
		MediaId       int    `json:"mediaId" binding:"number"`
		Answer        string `json:"answer" binding:"required,gte=4,lte=100"`
		AnswerImageId int    `json:"answerImageId" binding:"number"`
		LanguageId    int    `json:"languageId" binding:"required,number"`
	}
)

type (
	CreateDTO struct {
		Question      string
		MediaId       int
		Answer        string
		AnswerImageId int
		LanguageId    int
		AccountId     string
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}
)
