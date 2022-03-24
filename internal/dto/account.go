package dto

import "time"

type AccountVerification struct {
	AccountId string
	Code      string
	Verified  bool
	UpdatedAt time.Time
}
