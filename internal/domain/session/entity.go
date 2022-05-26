package session

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/answersuck/vault/pkg/strings"
)

type Session struct {
	Id        string    `json:"id"`
	AccountId string    `json:"-"`
	UserAgent string    `json:"userAgent"`
	IP        string    `json:"ip"`
	MaxAge    int       `json:"maxAge"`
	ExpiresAt int64     `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type fields struct {
	accountId  string
	userAgent  string
	ip         string
	expiration time.Duration
}

func newSession(f fields) (*Session, error) {
	// TODO: add useragent validation

	if _, err := netip.ParseAddr(f.ip); err != nil {
		return nil, fmt.Errorf("netip.ParseAddr: %w", err)
	}

	sid, err := strings.NewUnique(64)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Session{
		Id:        sid,
		AccountId: f.accountId,
		UserAgent: f.userAgent,
		IP:        f.ip,
		MaxAge:    int(f.expiration.Seconds()),
		ExpiresAt: now.Add(f.expiration).Unix(),
		CreatedAt: now,
	}, nil
}

type Device struct {
	UserAgent string
	IP        string
}
