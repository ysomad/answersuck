package apperr

import "errors"

const (
	MsgQuestionMediaNotFound = "question or answer media not found, upload it first"
)

var (
	QuestionNotFound = errors.New("question not found")
)
