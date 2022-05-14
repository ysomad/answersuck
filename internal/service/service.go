package service

import (
	"context"
	"time"

	"github.com/answersuck/vault/internal/domain"
)

type (
	SessionService interface {
		Create(ctx context.Context, accountId string, d domain.Device) (*domain.Session, error)
		GetById(ctx context.Context, sessionId string) (*domain.Session, error)
		GetAll(ctx context.Context, accountId string) ([]*domain.Session, error)
		Terminate(ctx context.Context, sessionId string) error
		TerminateAll(ctx context.Context, accountId string) error
		TerminateWithExcept(ctx context.Context, accountId, sessionId string) error
	}

	TokenManager interface {
		New(subject, audience string, expiration time.Duration) (string, error)
		Parse(token, audience string) (string, error)
	}
)
