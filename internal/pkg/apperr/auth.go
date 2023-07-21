package apperr

import "errors"

var (
	ErrNotAuthorized     = errors.New("invalid login or password")
	ErrAlreadyAuthorized = errors.New("already authorized, git gud")
)
