package apperr

import "errors"

const (
	MsgRoundQuestionNotFound = "round question not found"
)

var (
	RoundQuestionNotFound = errors.New(MsgRoundQuestionNotFound)
)
