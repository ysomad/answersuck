package language

import "errors"

var ErrNotFound = errors.New("language with given id not found")

type Language struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}
