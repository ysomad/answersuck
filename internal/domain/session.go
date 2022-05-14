package domain

import (
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/answersuck/vault/pkg/strings"
)

// Client errors
var (
	ErrSessionCannotBeTerminated = errors.New("current session cannot be terminated, use logout instead")
)

// System errors
var (
	ErrSessionContextNotFound     = errors.New("session not found in context")
	ErrSessionDeviceMismatch      = errors.New("device doesn't match with device of current session")
	ErrSessionAlreadyExist        = errors.New("session with given id already exist")
	ErrSessionForeignKeyViolation = errors.New("session cannot be created, account with given account id is not found")
	ErrSessionNotFound            = errors.New("session not found")
	ErrSessionNotDeleted          = errors.New("session has not been deleted")
)

type Session struct {
	Id        string    `json:"id"`
	AccountId string    `json:"accountId"`
	UserAgent string    `json:"userAgent"`
	IP        string    `json:"ip"`
	MaxAge    int       `json:"maxAge"`
	ExpiresAt int64     `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewSession(accountId, userAgent, ip string, expiration time.Duration) (*Session, error) {
	// TODO: add useragent validation

	if _, err := netip.ParseAddr(ip); err != nil {
		return nil, fmt.Errorf("netip.ParseAddr: %w", err)
	}

	sid, err := strings.NewUnique(64)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Session{
		Id:        sid,
		AccountId: accountId,
		UserAgent: userAgent,
		IP:        ip,
		MaxAge:    int(expiration.Seconds()),
		ExpiresAt: now.Add(expiration).Unix(),
		CreatedAt: now,
	}, nil
}
