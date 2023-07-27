package apperr

import "errors"

const (
	MsgMediaNotFound = "media not found"
)

var (
	MediaNotFound = errors.New(MsgMediaNotFound)
)
