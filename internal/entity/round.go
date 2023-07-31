package entity

// Maximum topics in one round.
const MaxRoundTopics = 10

type Round struct {
	ID       int32
	Name     string
	PackID   int32
	Position int16
}
