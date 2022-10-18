package domain

import (
	"errors"
	"time"
)

const EmailVerifCodeLen = 32

var (
	ErrAccountIDNotFound           = errors.New("account with given id not found")
	ErrEmailVerificationNotCreated = errors.New("error occured on email verification create")
)

type EmailVerification struct {
	AccountID string
	Code      string
	ExpiresAt time.Time
}

func NewEmailVerification(accountID, code string, expiresIn time.Duration) EmailVerification {
	return EmailVerification{
		AccountID: accountID,
		Code:      code,
		ExpiresAt: time.Now().Add(expiresIn),
	}
}
