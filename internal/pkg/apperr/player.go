package apperr

import "errors"

var (
	ErrPlayerNotFound      = errors.New("player not found")
	ErrPlayerAlreadyExists = errors.New("player with provided email or nickname already exists")
)
