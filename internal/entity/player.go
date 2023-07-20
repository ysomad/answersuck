package entity

import "time"

type Player struct {
	Nickname      string
	Email         string
	DisplayName   string
	EmailVerified bool
	PasswordHash  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
