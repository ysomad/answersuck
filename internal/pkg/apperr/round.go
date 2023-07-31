package apperr

import "errors"

const (
	MsgRoundNotFound       = "round not found"
	MsgRoundTopicNotAdded  = "amount of topic in round exceeded"
	MsgTopicAlreadyInRound = "topic already added to round"
)

var (
	RoundNotFound       = errors.New(MsgRoundNotFound)
	RoundTopicNotAdded  = errors.New(MsgRoundTopicNotAdded)
	TopicAlreadyInRound = errors.New(MsgTopicAlreadyInRound)
)
