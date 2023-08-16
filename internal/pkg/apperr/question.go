package apperr

import "errors"

const (
	MsgQuestionMediaNotFound = "question or answer media not found, upload it first"
	MsgQuestionNotFound      = "question not found"
)

var (
	QuestionNotFound = errors.New(MsgQuestionNotFound)
)
