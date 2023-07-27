package apperr

import "errors"

var (
	QuestionNotFound      = errors.New("question not found")
	QuestionMediaNotExist = errors.New("question media doesn't exist")
	AnswerMediaNotExist   = errors.New("answer media doesn't exist")
)
