package domain

import (
	"errors"
	"time"

	"github.com/ysomad/answersuck/cryptostr"
)

const (
	passwordTokenLen = 128
)

var (
	ErrPasswordTokenExpired = errors.New("token expired")
)

type PasswordToken struct {
	AccountID string
	Token     string
	ExpiresAt time.Time
}

func GenPasswordToken() (string, error) {
	return cryptostr.RandomWithAlphabetDigits(passwordTokenLen)
}
