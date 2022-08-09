package question

import "time"

type CreateDTO struct {
	Text       string
	AnswerId   uint32
	MediaId    *string
	AccountId  string
	LanguageId uint8
	CreatedAt  time.Time
}
