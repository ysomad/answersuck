package entity

import (
	"errors"
	"time"
)

type QuestionTranfserType int8

const (
	QTranfserTypeUnspecified QuestionTranfserType = iota
	QTranfserTypeBefore
	QTranfserTypeAfter
	QTranfserTypeNever
)

func (t QuestionTranfserType) valid() bool {
	switch t {
	case QTranfserTypeAfter, QTranfserTypeNever, QTranfserTypeBefore:
		return true
	}

	return false
}

type RoundQuestion struct {
	ID           int32
	QuestionID   int32
	TopicID      int32
	RoundID      int32
	Type         QuestionType
	Cost         int32
	AnswerTime   time.Duration
	HostComment  string
	SecretTopic  string
	SecretCost   int32
	Keepable     bool
	TransferType QuestionTranfserType
}

var (
	ErrInvalidSecretQuestion      = errors.New("invalid secret question")
	ErrInvalidSuperSecretQuestion = errors.New("invalid super secret question")
)

func (q *RoundQuestion) Validate() error {
	switch q.Type {
	case QTypeSecret:
		if q.SecretCost == 0 || q.SecretTopic == "" || q.Keepable {
			return ErrInvalidSecretQuestion
		}
	case QTypeSuperSecret:
		if q.SecretCost == 0 || q.SecretTopic == "" || !q.TransferType.valid() {
			return ErrInvalidSuperSecretQuestion
		}
	}

	return nil
}
