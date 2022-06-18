package account

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/answersuck/vault/pkg/strings"
)

var (
	ErrAlreadyExist              = errors.New("account with given email or username already exist")
	ErrNotDeleted                = errors.New("account has not been deleted")
	ErrAlreadyArchived           = errors.New("account already archived or not found")
	ErrIncorrectCredentials      = errors.New("incorrect login or password")
	ErrAlreadyVerified           = errors.New("current email already verified or verification code is expired")
	ErrForbiddenNickname         = errors.New("nickname contains forbidden words")
	ErrNotFound                  = errors.New("account not found")
	ErrEmptyVerificationCode     = errors.New("empty account verification code")
	ErrEmptyPasswordResetToken   = errors.New("empty password reset token")
	ErrPasswordTokenNotFound     = errors.New("account password reset token not found or expired")
	ErrIncorrectPassword         = errors.New("incorrect password")
	ErrPasswordResetTokenExpired = errors.New("password reset token is expired")
	ErrVerificationNotFound      = errors.New("account verification not found")
	ErrPasswordNotSet            = errors.New("account password is not set")
	ErrPasswordTokenAlreadyExist = errors.New("account password reset token already exist")
	ErrNotEnoughRights           = errors.New("not enough rights to perform the operation, verify account first")
)

type Account struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"-"`
	Verified  bool      `json:"verified"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Account) setPassword(p string) error {
	h, err := a.hashPassword(p)
	if err != nil {
		return err
	}

	a.Password = h

	return nil
}

func (a *Account) ComparePasswords(p string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(p)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrIncorrectPassword)
	}

	return nil
}

func (a *Account) hashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), 11)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (a *Account) generateVerifCode(length int) (string, error) {
	return strings.NewUnique(length)
}

type PasswordToken struct {
	AccountId string
	Token     string
	CreatedAt time.Time
}

// checkExpiration returns error if token is expired
func (t PasswordToken) checkExpiration(exp time.Duration) error {
	expiresAt := t.CreatedAt.Add(exp)

	if time.Now().After(expiresAt) {
		return ErrPasswordResetTokenExpired
	}

	return nil
}

type Verification struct {
	Email    string
	Code     string
	Verified bool
}
