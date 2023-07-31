package apperr

import "errors"

const (
	MsgRoundNotFound = "round not found"

	MsgRoundTopicNotAdded      = "amount of topic in round exceeded"
	MsgRoundTopicAlreadyExists = "topic already added to round"
	MsgRoundTopicNotDeleted    = "round or topic not found"
)

var (
	RoundNotFound = errors.New(MsgRoundNotFound)

	RoundTopicNotAdded      = errors.New(MsgRoundTopicNotAdded)
	RoundTopicAlreadyExists = errors.New(MsgRoundTopicAlreadyExists)
	RoundTopicNotDeleted    = errors.New(MsgRoundTopicNotDeleted)
)
