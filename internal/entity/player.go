package entity

import (
	"net/mail"
	"time"
)

type LoginType int8

const (
	LoginTypeNickname LoginType = iota + 1
	LoginTypeEmail
)

func NewLoginType(login string) LoginType {
	if _, err := mail.ParseAddress(login); err == nil {
		return LoginTypeEmail
	}

	return LoginTypeNickname
}

type Player struct {
	Nickname      string
	Email         string
	DisplayName   string
	EmailVerified bool
	PasswordHash  string
	CreatedAt     time.Time
	UpdateTime    time.Time
}
