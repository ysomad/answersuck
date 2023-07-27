package apperr

import "errors"

const (
	MsgTagAlreadyExists = "tag with provided name already exists"
)

var (
	TagAlreadyExists = errors.New(MsgTagAlreadyExists)
)
