package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/quizlyfun/quizly-backend/pkg/utils"
)

var (
	// Client errors
	ErrAccountAlreadyExist         = errors.New("account with given email or username already exist")
	ErrAccountIncorrectCredentials = errors.New("provided credentials are not correct")
	ErrAccountEmptyEmailOrUsername = errors.New("email or username should be provided")

	// System errors
	ErrAccountNotFound             = errors.New("account not found")
	ErrAccountIncorrectPassword    = errors.New("incorrect password")
	ErrAccountPasswordNotGenerated = errors.New("password hash generation error")
	ErrAccountNotArchived          = errors.New("account cannot be archived")
	ErrAccountContextNotFound      = errors.New("account not found in context")
	ErrAccountContextMismatch      = errors.New("account id from context is not the same as account id from url parameter")
)

type Account struct {
	Id           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	Verified     bool      `json:"verified"`
	Archived     bool      `json:"archived"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// GeneratePasswordHash generates hash from password and sets it to PasswordHash
func (a *Account) GeneratePasswordHash() error {
	b, err := bcrypt.GenerateFromPassword([]byte(a.Password), 11)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", ErrAccountPasswordNotGenerated)
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
	a.Password = utils.RandomSpecialString(16)
}
