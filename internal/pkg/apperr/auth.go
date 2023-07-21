package apperr

import "errors"

var (
	ErrNotAuthorized = errors.New("invalid login or password")
)
