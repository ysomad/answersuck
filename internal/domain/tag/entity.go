package tag

import "errors"

var (
	ErrLanguageIdNotFound = errors.New("language id not found")
	ErrEmptyTagList       = errors.New("empty tag list")
)

type Tag struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	LanguageId int    `json:"language_id"`
}
