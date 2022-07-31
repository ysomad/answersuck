package account

import (
	"errors"
	"time"
)

var (
	ErrAlreadyExist      = errors.New("account with given email or nickname already exist")
	ErrNotDeleted        = errors.New("account has not been deleted")
	ErrAlreadyArchived   = errors.New("account already archived or not found")
	ErrAlreadyVerified   = errors.New("current email already verified or verification code is expired")
	ErrForbiddenNickname = errors.New("nickname contains forbidden words")
	ErrNotFound          = errors.New("account not found")
	ErrNotEnoughRights   = errors.New("account must be verified to perform this operation")

	ErrEmptyVerificationCode = errors.New("empty account verification code")
	ErrVerificationNotFound  = errors.New("account verification not found")

	ErrEmptyPasswordToken        = errors.New("empty password reset token")
	ErrPasswordTokenNotFound     = errors.New("account password reset token not found or expired")
	ErrPasswordTokenExpired      = errors.New("password reset token is expired")
	ErrPasswordTokenAlreadyExist = errors.New("account password reset token already exist")

	ErrPasswordNotSet = errors.New("account password is not set")
)

type Account struct {
	Id        string
	Email     string
	Nickname  string
	Password  string
	Verified  bool
	Archived  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	VerifCodeLen     = 64
	PasswordTokenLen = 64
)

type PasswordToken struct {
	AccountId string
	Token     string
	CreatedAt time.Time
}

func (t PasswordToken) expired(exp time.Duration) bool {
	return time.Now().After(t.CreatedAt.Add(exp))
}

type Verification struct {
	Email    string
	Code     string
	Verified bool
}
