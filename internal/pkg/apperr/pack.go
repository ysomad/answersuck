package apperr

import "errors"

const (
	MsgPackCoverNotFound = "pack cover not found, upload media first"
	MsgPackNotFound      = "pack not found"
	MsgPackNotAuthor     = "current user is not an author of the pack"
)

var (
	PackNotFound  = errors.New(MsgPackNotFound)
	PackNotAuthor = errors.New(MsgPackNotAuthor)
)
