package question

import "errors"

var (
	ErrForeignKeyViolation = errors.New("provided answer, author account, media or language are not exist")
	ErrNotFound            = errors.New("question with provided id not found")
)
