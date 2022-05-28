package topic

import (
	"time"
)

type Topic struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	LanguageId int       `json:"languageId"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
