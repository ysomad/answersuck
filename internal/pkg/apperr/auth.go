package apperr

import "errors"

const (
	MsgUnauthorized       = "unauthorized"
	MsgAuthorized         = "already authorized"
	MsgInvalidCredentials = "invalid login or password"
)

var (
	ErrInvalidCredentials = errors.New(MsgInvalidCredentials)
)
