package dto

import "time"

type UpdateEmailArgs struct {
	AccountID, NewEmail, PlainPassword string
}

type EmailVerifSaveArgs struct {
	Code      string
	ExpiresAt time.Time
}
