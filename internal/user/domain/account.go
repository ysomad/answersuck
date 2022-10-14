package domain

import (
	"errors"
	"time"

	"github.com/ysomad/answersuck/internal/user/service/dto"
)

var (
	ErrAccountAlreadyExist = errors.New("account with given email or username already exist")
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
