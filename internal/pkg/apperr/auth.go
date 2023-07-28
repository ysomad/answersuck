package apperr

import "errors"

const (
	MsgUnauthorized       = "unauthorized"
	MsgAuthorized         = "already authorized"
	MsgInvalidCredentials = "invalid login or password"
)

var (
	InvalidCredentials = errors.New(MsgInvalidCredentials)
	Unauthorized       = errors.New(MsgUnauthorized)
)
