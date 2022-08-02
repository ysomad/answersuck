package session

import (
	"errors"
	"time"

	"github.com/answersuck/host/internal/pkg/strings"
)

var (
	ErrCannotBeTerminated = errors.New("current session cannot be terminated, use logout instead")
	ErrDeviceMismatch     = errors.New("device doesn't match with device of current session")
	ErrAlreadyExist       = errors.New("session with given id already exist")
	ErrAccountNotFound    = errors.New("session cannot be created, account with given account id is not found")
	ErrNotFound           = errors.New("session not found")
	ErrNotDeleted         = errors.New("session has not been deleted")
	ErrExpired            = errors.New("session expired")
)

const SessionIdLen = 64

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
	// TODO: add ip validation

	sid, err := strings.NewUnique(SessionIdLen)
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

func (s *Session) SameDevice(d Device) bool {
	if s.IP != d.IP || s.UserAgent != d.UserAgent {
		return false
	}
	return true
}

type Device struct {
	UserAgent string
	IP        string
}

type WithAccountDetails struct {
	Session
	AccountVerified bool
}
