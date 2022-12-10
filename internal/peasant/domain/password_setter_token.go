package domain

import (
	"errors"
)

var (
	ErrPasswordSetterTokenExpired = errors.New("password setter token expired")
)

// PasswordSetterToken token must be used when user forgot his password,
// using the token its possible to update the password,
// must be created only via constructor.
type PasswordSetterToken string

func (t PasswordSetterToken) String() string { return string(t) }
