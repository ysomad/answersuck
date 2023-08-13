package apperr

import "errors"

const (
	MsgTopicNotFound = "topic not found"
)

var (
	TopicNotFound = errors.New(MsgTopicNotFound)
)
