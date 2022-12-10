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
)

type Account struct {
	ID            string
	Email         string
	Username      string
	EmailVerified bool
	Archived      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
