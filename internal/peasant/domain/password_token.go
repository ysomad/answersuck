package domain

import (
	"time"

	"github.com/ysomad/answersuck/cryptostr"
)

const (
	passwordTokenLen = 64
)

type PasswordToken struct {
	AccountID string
	Token     string
	ExpiresAt time.Time
}

func GenPasswordToken() (string, error) { return cryptostr.RandomBase64(passwordTokenLen) }
