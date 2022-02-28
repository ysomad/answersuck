package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/quizlyfun/quizlyfun-backend/pkg/utils"
)

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionExpired         = errors.New("session expired")
	ErrSessionNotCreated      = errors.New("error occured during session creation")
	ErrSessionNotTerminated   = errors.New("current session cannot be terminated, use logout instead")
	ErrSessionContextNotFound = errors.New("session not found in context")
	ErrSessionDeviceMismatch  = errors.New("device doesn't match with device of current session")
)

type Session struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"accountId" bson:"accountId"`
	Device    Device    `json:"device" bson:"device"`
	TTL       int       `json:"ttl" bson:"ttl"`
	ExpiresAt int64     `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type Device struct {
	UserAgent string
	IP        string
}

func NewSession(accountID string, d Device, ttl time.Duration) (Session, error) {
	id, err := utils.UniqueString(32)
	if err != nil {
		return Session{}, fmt.Errorf("utils.UniqueString: %w", ErrSessionNotCreated)
	}

	now := time.Now()

	return Session{
		ID:        id,
		AccountID: accountID,
		Device:    d,
		TTL:       int(ttl.Seconds()),
		ExpiresAt: now.Add(ttl).Unix(),
		CreatedAt: now,
	}, nil
}
