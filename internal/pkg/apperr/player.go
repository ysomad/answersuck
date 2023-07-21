package apperr

import "errors"

var (
	ErrPlayerNotFound     = errors.New("player not found")
	ErrPlayerAlreadyExist = errors.New("player with provided email or nickname already exist")
)
