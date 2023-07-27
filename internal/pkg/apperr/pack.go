package apperr

import "errors"

const (
	MsgPackCoverNotFound = "pack cover not found, upload media first"
	MsgPackNotFound      = "pack not found"
)

var (
	PackNotFound = errors.New(MsgPackNotFound)
)
