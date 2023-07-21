package entity

import "time"

type LoginType int8

const (
	LoginTypeNickname LoginType = iota + 1
	LoginTypeEmail
)

type Player struct {
	Nickname      string
	Email         string
	DisplayName   string
	EmailVerified bool
	PasswordHash  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
