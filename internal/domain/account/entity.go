package account

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/answersuck/vault/pkg/dicebear"
	"github.com/answersuck/vault/pkg/strings"
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
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", ErrIncorrectPassword)
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

type PasswordResetToken struct {
	AccountId string
	Token     string
	CreatedAt time.Time
}

// CheckExpiration returns error if token is expired
func (t PasswordResetToken) CheckExpiration(exp time.Duration) error {
	expiresAt := t.CreatedAt.Add(exp)

	if time.Now().After(expiresAt) {
		return ErrPasswordResetTokenExpired
	}

	return nil
}
