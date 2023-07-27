package apperr

import "errors"

const (
	MsgPlayerNotVerified = "player not verified"
)

var (
	PlayerNotFound      = errors.New("player not found")
	PlayerAlreadyExists = errors.New("player with provided email or nickname already exists")
)
