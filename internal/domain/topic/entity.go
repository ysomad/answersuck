package topic

import (
	"time"
)

type Topic struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	LanguageId int       `json:"language_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
