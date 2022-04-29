package domain

import "errors"

var ErrLanguageNotFound = errors.New("language with given id not found")

type Language struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
