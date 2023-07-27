package apperr

import "errors"

var (
	QuestionMediaNotExist = errors.New("question media doesn't exist")
	AnswerMediaNotExist   = errors.New("answer media doesn't exist")
)
