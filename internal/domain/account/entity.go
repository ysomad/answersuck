package account

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/answersuck/vault/pkg/strings"
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

func (a *Account) ComparePasswords(p string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(p)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrIncorrectPassword)
	}

	return nil
}

// setPassword hashes given password and sets it to password field
func (a *Account) setPassword(p string) error {
	ph, err := generatePasswordHash(p)
	if err != nil {
		return fmt.Errorf("generatePasswordHash: %w", err)
	}

	a.Password = ph

	return nil
}

func generatePasswordHash(p string) (string, error) {
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
