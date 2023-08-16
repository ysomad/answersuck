package entity

import (
	"errors"
	"time"
)

type QuestionType int8

const (
	QTypeStandard QuestionType = iota + 1
	QTypeSafe
	QTypeSecret
	QTypeSuperSecret
	QTypeAuction
)

type QuestionTransferType int8

const (
	QTransferTypeUnspecified QuestionTransferType = iota
	QTransferTypeBefore
	QTransferTypeAfter
	QTransferTypeNever
)

func (t QuestionTransferType) valid() bool {
	switch t {
	case QTransferTypeAfter, QTransferTypeNever, QTransferTypeBefore:
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
	TransferType QuestionTransferType
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

type RoundQuestionDetailed struct {
	RoundQuestion

	Question         string
	QuestionMediaURL string

	AnswerID       int32
	Answer         string
	AnswerMediaURL string
}
