package auth

import "errors"

var (
	ErrAlreadyLoggedIn          = errors.New("already logged in")
	ErrIncorrectCredentials     = errors.New("incorrect login or password")
	ErrIncorrectAccountPassword = errors.New("incorrect account password")
)
