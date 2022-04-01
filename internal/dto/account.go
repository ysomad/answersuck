package dto

import (
	"time"
)

type AccountVerification struct {
	Email    string
	Code     string
	Verified bool
}

type AccountPasswordResetToken struct {
	AccountId string
	Token     string
	CreatedAt time.Time
}

type AccountUpdatePassword struct {
	AccountId    string
	Token        string
	PasswordHash string
	UpdatedAt    time.Time
}
