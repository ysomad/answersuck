package apperr

import "errors"

var (
	ErrQestionMediaNotExist = errors.New("question media doesn't exist")
	ErrAnswerMediaNotExist  = errors.New("answer media doesn't exist")
)
