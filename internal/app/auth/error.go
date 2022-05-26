package auth

import "errors"

var (
	ErrAlreadyLoggedIn = errors.New("already logged in")
)
