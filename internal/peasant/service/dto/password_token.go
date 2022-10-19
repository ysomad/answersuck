package dto

import "time"

type CreatePasswordTokenArgs struct {
	EmailOrUsername string
	Token           string
	ExpiresAt       time.Time
}

func NewCreatePasswordTokenArgs(emailOrUsername, token string, expiresIn time.Duration) CreatePasswordTokenArgs {
	return CreatePasswordTokenArgs{
		EmailOrUsername: emailOrUsername,
		Token:           token,
		ExpiresAt:       time.Now().Add(expiresIn),
	}
}
