package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"
)

var (
	errInvalidUserID  = errors.New("session: invalid user id")
	errEmptyUserAgent = errors.New("session: empty user agent")
)

type Store interface {
	// Save creates new session, returns session id.
	Save(context.Context, *Session) error

	// Get returns session by its id.
	Get(context.Context, string) (*Session, error)

	// Delete deletes session by its id.
	Delete(context.Context, string) error
}

type Manager struct {
	store    Store
	lifetime time.Duration
}

func NewManager(s Store, lifetime time.Duration) *Manager {
	return &Manager{
		store:    s,
		lifetime: lifetime,
	}
}

func (m *Manager) Create(ctx context.Context, p User) (*Session, error) {
	if p.ID == "" {
		return nil, errInvalidUserID
	}

	if p.UserAgent == "" {
		return nil, errEmptyUserAgent
	}

	sid, err := m.generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sid,
		User:      p,
		ExpiresAt: time.Now().Add(m.lifetime),
	}

	if err := m.store.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (m *Manager) Get(ctx context.Context, sid string) (*Session, error) {
	return m.store.Get(ctx, sid)
}

func (m *Manager) Delete(ctx context.Context, sid string) error {
	return m.store.Delete(ctx, sid)
}

func (m *Manager) generateSessionID() (string, error) {
	b := make([]byte, 64)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
