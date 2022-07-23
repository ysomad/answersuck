package auth

import "errors"

var (
	ErrAlreadyLoggedIn      = errors.New("already logged in")
	ErrIncorrectPassword    = errors.New("incorrect account password")
	ErrAccountNotFound      = errors.New("accout not found")
	ErrIncorrectCredentials = errors.New("incorrect login or password")
)
