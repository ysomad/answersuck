package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/answersuck/answersuck-backend/pkg/strings"
)

// Client errors
var (
	ErrSessionAlreadyLoggedIn = errors.New("already logged in, please logout before logging in")
)

// System errors
var (
	ErrSessionNotCreated      = errors.New("error occurred during session creation")
	ErrSessionContextNotFound = errors.New("session not found in context")
	ErrSessionDeviceMismatch  = errors.New("device doesn't match with device of current session")
)

type Session struct {
	Id         string        `json:"id"`
	AccountId  string        `json:"accountId"`
	Device     Device        `json:"device"`
	MaxAge     int           `json:"maxAge"`
	Expiration time.Duration `json:"expiration"`
	ExpiresAt  int64         `json:"expiresAt"`
	CreatedAt  time.Time     `json:"createdAt"`
}

func NewSession(accountId string, d Device, expiration time.Duration) (*Session, error) {
	if err := d.Validate(); err != nil {
		return nil, fmt.Errorf("d.Validate: %w", err)
	}

	id, err := strings.NewUnique(32)
	if err != nil {
		return nil, fmt.Errorf("utils.UniqueString: %w", ErrSessionNotCreated)
	}

	now := time.Now()

	return &Session{
		Id:         id,
		AccountId:  accountId,
		Device:     d,
		MaxAge:     int(expiration.Seconds()),
		Expiration: expiration,
		ExpiresAt:  now.Add(expiration).Unix(),
		CreatedAt:  now,
	}, nil
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

var (
	ErrDeviceInvalidUserAgent = errors.New("invalid user agent")
	ErrDeviceInvalidIP        = errors.New("invalid ip address")
)

type Device struct {
	UserAgent string `json:"ua"`
	IP        string `json:"ip"`
}

func (d Device) validateIP() error {
	if _, err := netip.ParseAddr(d.IP); err != nil {
		return ErrDeviceInvalidIP
	}

	return nil
}

func (d Device) Validate() error {
	// TODO: add validation for User-Agent
	if d.UserAgent == "" {
		return ErrDeviceInvalidUserAgent
	}

	if err := d.validateIP(); err != nil {
		return fmt.Errorf("d.validateIP: %w", err)
	}

	return nil
}
