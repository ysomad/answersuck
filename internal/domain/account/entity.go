package account

import (
	"errors"
	"time"

	"github.com/answersuck/vault/pkg/strings"
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

func (a *Account) generateVerifCode(length int) (string, error) {
	return strings.NewUnique(length)
}

const (
	verifCodeLen = 64
	pwdTokenLen  = 64
)

type PasswordToken struct {
	AccountId string
	Token     string
	CreatedAt time.Time
}

// checkExpiration returns error if token is expired
func (t PasswordToken) checkExpiration(exp time.Duration) error {
	if time.Now().After(t.CreatedAt.Add(exp)) {
		return ErrPasswordTokenExpired
	}

	return nil
}

type Verification struct {
	Email    string
	Code     string
	Verified bool
}
