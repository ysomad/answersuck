package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/answersuck/answersuck-backend/pkg/dicebear"
	"github.com/answersuck/answersuck-backend/pkg/strings"
)

// Client errors
var (
	ErrAccountAlreadyExist            = errors.New("account with given email or username already exist")
	ErrAccountIncorrectCredentials    = errors.New("incorrect login or password")
	ErrAccountAlreadyVerified         = errors.New("current email already verified or verification code is expired")
	ErrAccountForbiddenUsername       = errors.New("username contains forbidden words")
	ErrAccountNotFound                = errors.New("account not found")
	ErrAccountEmptyVerificationCode   = errors.New("empty account verification code")
	ErrAccountEmptyResetPasswordToken = errors.New("empty reset password token")
)

// System errors
var (
	ErrAccountIncorrectPassword         = errors.New("incorrect password")
	ErrAccountContextNotFound           = errors.New("account not found in context")
	ErrAccountResetPasswordTokenExpired = errors.New("password reset token is expired")
)

type Account struct {
	Id               string    `json:"id"`
	Email            string    `json:"email"`
	Username         string    `json:"username"`
	Password         string    `json:"-"`
	PasswordHash     string    `json:"-"`
	Verified         bool      `json:"verified"`
	VerificationCode string    `json:"-"`
	Archived         bool      `json:"archived"`
	AvatarURL        string    `json:"avatarUrl"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// GeneratePasswordHash generates hash from password and sets it to PasswordHash
func (a *Account) GeneratePasswordHash() error {
	b, err := bcrypt.GenerateFromPassword([]byte(a.Password), 11)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	a.PasswordHash = string(b)

	return nil
}

func (a *Account) CompareHashAndPassword() error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(a.Password)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrAccountIncorrectPassword)
	}

	return nil
}

// SetDiceBearAvatar sets dicebear identicon url from username to account AvatarURL
func (a *Account) SetDiceBearAvatar() {
	a.AvatarURL = dicebear.URL(a.Username)
}

func (a *Account) GenerateVerificationCode(length int) error {
	code, err := strings.NewUnique(length)
	if err != nil {
		return fmt.Errorf("strings.NewUnique: %w", err)
	}

	a.VerificationCode = code

	return nil
}
