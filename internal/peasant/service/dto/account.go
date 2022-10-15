package dto

import "time"

type AccountSaveArgs struct {
	Email           string
	Username        string
	EncodedPassword string
}

type EmailVerifSaveArgs struct {
	Code      string
	ExpiresAt time.Time
}

type AccountCreateArgs struct {
	Email, Username, PlainPassword string
}
