package apperr

import "errors"

const (
	MsgRoundNotFound = "round not found"
)

var (
	RoundNotFound = errors.New(MsgRoundNotFound)
)
