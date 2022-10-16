package domain

import (
	"errors"
	"time"
)

var (
	ErrUsernameTaken     = errors.New("username already in use")
	ErrEmailTaken        = errors.New("email already in use")
	ErrAccountNotFound   = errors.New("account not found")
	ErrIncorrectPassword = errors.New("incorrect password")

	ErrEmailNotVerified = errors.New("email already verified or code expired")
)

type Account struct {
	ID            string
	Username      string
	Email         string
	EmailVerified bool
	Archived      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
