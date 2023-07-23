package apperr

import "errors"

const (
	MsgPlayerNotVerified = "player not verified"
)

var (
	ErrPlayerNotFound      = errors.New("player not found")
	ErrPlayerAlreadyExists = errors.New("player with provided email or nickname already exists")
)
