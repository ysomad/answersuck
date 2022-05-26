package email

import "errors"

var (
	ErrEmptyTo      = errors.New("empty to")
	ErrEmptySubject = errors.New("empty subject")
	ErrEmptyMessage = errors.New("empty message")
)
