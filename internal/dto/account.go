package dto

import "time"

type AccountVerify struct {
	AccountId string
	Code      string
	Verified  bool
	UpdatedAt time.Time
}

type AccountArchive struct {
	AccountId string
	Archived  bool
	UpdatedAt time.Time
}
