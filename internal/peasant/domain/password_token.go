package domain

import (
	"errors"

	"github.com/ysomad/answersuck/jwt"
)

var (
	ErrPasswordTokenExpired = errors.New("token expired")
)

// PasswordToken token must be used when user forgot his password,
// using the token its possible to update the password,
// must be created only via constructor.
type PasswordToken string

func NewPasswordToken(b jwt.Basic) PasswordToken {
	return PasswordToken(b)
}

func (t PasswordToken) String() string { return string(t) }
