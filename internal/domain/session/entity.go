package session

import (
	"time"

	"github.com/answersuck/vault/pkg/strings"
)

type Session struct {
	Id        string    `json:"id"`
	AccountId string    `json:"-"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	MaxAge    int       `json:"-"`
	ExpiresAt int64     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func newSession(accountId, userAgent, ip string, exp time.Duration) (*Session, error) {
	// TODO: add useragent validation

	// if _, err := netip.ParseAddr(ip); err != nil {
	// 	return nil, fmt.Errorf("netip.ParseAddr: %w", err)
	// }

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
		MaxAge:    int(exp.Seconds()),
		ExpiresAt: now.Add(exp).Unix(),
		CreatedAt: now,
	}, nil
}

func (s *Session) Expired() bool { return time.Now().Unix() > s.ExpiresAt }

func (s *Session) SameDevice(ip, ua string) bool {
	if s.IP != ip || s.UserAgent != ua {
		return false
	}
	return true
}

type Device struct {
	UserAgent string
	IP        string
}
