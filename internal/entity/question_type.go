package entity

type QuestionType int8

const (
	QTypeStandard QuestionType = iota + 1
	QTypeSafe
	QTypeSecret
	QTypeSuperSecret
	QTypeAuction
)
