package session

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/answersuck/vault/pkg/strings"
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

func newSession(accountId, userAgent, ip string, expiration time.Duration) (*Session, error) {
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

type Device struct {
	UserAgent string
	IP        string
}
