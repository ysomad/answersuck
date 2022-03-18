package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/quizlyfun/quizly-backend/pkg/utils"
)

var (
	// Client errors
	ErrSessionAlreadyLoggedIn = errors.New("already logged in, please logout before logging in")

	// System errors
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionExpired         = errors.New("session expired")
	ErrSessionNotCreated      = errors.New("error occured during session creation")
	ErrSessionNotTerminated   = errors.New("current session cannot be terminated, use logout instead")
	ErrSessionContextNotFound = errors.New("session not found in context")
	ErrSessionDeviceMismatch  = errors.New("device doesn't match with device of current session")
)

type Session struct {
	Id         string        `json:"id"`
	AccountId  string        `json:"accountId"`
	Device     Device        `json:"device"`
	TTL        int           `json:"ttl"`
	Expiration time.Duration `json:"expiration"`
	ExpiresAt  int64         `json:"expiresAt"`
	CreatedAt  time.Time     `json:"createdAt"`
}

func NewSession(accountId string, d Device, expiration time.Duration) (Session, error) {
	if err := d.Validate(); err != nil {
		return Session{}, fmt.Errorf("d.Validate: %w", err)
	}

	id, err := utils.UniqueString(32)
	if err != nil {
		return Session{}, fmt.Errorf("utils.UniqueString: %w", ErrSessionNotCreated)
	}

	now := time.Now()

	return Session{
		Id:         id,
		AccountId:  accountId,
		Device:     d,
		TTL:        int(expiration.Seconds()),
		Expiration: expiration,
		ExpiresAt:  now.Add(expiration).Unix(),
		CreatedAt:  now,
	}, nil
}

func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
