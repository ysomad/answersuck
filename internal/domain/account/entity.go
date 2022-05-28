package account

import (
	"fmt"
	"time"

	"github.com/answersuck/vault/pkg/strings"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id           string    `json:"id"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	Verified     bool      `json:"verified"`
	Archived     bool      `json:"archived"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (a *Account) CompareHashAndPassword() error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(a.Password)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrIncorrectPassword)
	}

	return nil
}

// generatePasswordHash generates hash from password and sets it to PasswordHash
func (a *Account) generatePasswordHash() error {
	b, err := bcrypt.GenerateFromPassword([]byte(a.Password), 11)
	if err != nil {
		return fmt.Errorf("bcrypt.generateFromPassword: %w", err)
	}

	a.PasswordHash = string(b)

	return nil
}

func (a *Account) generateVerificationCode(length int) (string, error) {
	return strings.NewUnique(length)
}

type PasswordResetToken struct {
	AccountId string
	Token     string
	CreatedAt time.Time
}

// checkExpiration returns error if token is expired
func (t PasswordResetToken) checkExpiration(exp time.Duration) error {
	expiresAt := t.CreatedAt.Add(exp)

	if time.Now().After(expiresAt) {
		return ErrPasswordResetTokenExpired
	}

	return nil
}
