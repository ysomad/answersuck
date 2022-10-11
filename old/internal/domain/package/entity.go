package packages

import "time"

type Package struct {
	Id          uint32
	Name        string
	Description string
	Published   bool
	AccountId   string
	LanguageId  uint8
	Tags        []uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
