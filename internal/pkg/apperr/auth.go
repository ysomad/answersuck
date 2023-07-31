package apperr

import "errors"

const (
	MsgUnauthorized         = "unauthorized"
	MsgAuthorized           = "already authorized"
	MsgInvalidCredentials   = "invalid login or password"
	MsgInvalidXRealIPHeader = "invalid or empty X-Real-IP header"
)

var (
	InvalidCredentials = errors.New(MsgInvalidCredentials)
	Unauthorized       = errors.New(MsgUnauthorized)
)
