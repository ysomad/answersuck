package dto

import "time"

type UpdateEmailArgs struct {
	AccountID, NewEmail, PlainPassword string
}

type EmailVerifCreateArgs struct {
	Code      string
	ExpiresAt time.Time
}

func NewEmailVerifCreateArgs(code string, expiresIn time.Duration) EmailVerifCreateArgs {
	return EmailVerifCreateArgs{
		Code:      code,
		ExpiresAt: time.Now().Add(expiresIn),
	}
}
