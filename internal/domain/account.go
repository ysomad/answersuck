package domain

import (
	"errors"
	"fmt"
	"github.com/quizlyfun/quizly-backend/pkg/dicebear"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/quizlyfun/quizly-backend/pkg/strings"
)

// Client errors
var (
	ErrAccountAlreadyExist         = errors.New("account with given email or username already exist")
	ErrAccountIncorrectCredentials = errors.New("incorrect login or password")
	ErrAccountAlreadyVerified      = errors.New("current email already verified or verification code is expired")
	ErrAccountForbiddenUsername    = errors.New("username contains forbidden words")
)

// System errors
var (
	ErrAccountNotFound              = errors.New("account not found")
	ErrAccountIncorrectPassword     = errors.New("incorrect password")
	ErrAccountContextNotFound       = errors.New("account not found in context")
	ErrAccountEmptyVerificationCode = errors.New("empty account verification code")
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

func (a *Account) RandomPassword() {
	a.Password = strings.NewSpecialRandom(16)
}

// DiceBearAvatar sets dicebear identicon url from username to account AvatarURL
func (a *Account) DiceBearAvatar() {
	a.AvatarURL = dicebear.URL(a.Username)
}

func (a *Account) GenerateVerificationCode() error {
	code, err := strings.NewUnique(32)
	if err != nil {
		return fmt.Errorf("strings.NewUnique: %w", err)
	}

	a.VerificationCode = code

	return nil
}
