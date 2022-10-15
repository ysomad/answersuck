package domain

import (
	"errors"
	"time"

	"github.com/ysomad/answersuck/internal/user/service/dto"
)

var (
	ErrUsernameTaken = errors.New("username already in use")
	ErrEmailTaken    = errors.New("email already in use")
)

type Account struct {
	ID              string
	Username        string
	Email           string
	EmailVerified   bool
	EncodedPassword string
	Archived        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// WithSaveArgs sets all fields from dto.AccountSaveArgs to Account.
func (a *Account) WithSaveArgs(args dto.AccountSaveArgs) *Account {
	a.Email = args.Email
	a.Username = args.Username
	a.EncodedPassword = args.EncodedPassword
	return a
}
