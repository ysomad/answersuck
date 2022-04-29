package service

import (
	"context"
	"time"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
)

type (
	sessionService interface {
		Create(ctx context.Context, aid string, d dto.Device) (*domain.Session, error)
		GetById(ctx context.Context, sid string) (*domain.Session, error)
		GetAll(ctx context.Context, aid string) ([]*domain.Session, error)
		Terminate(ctx context.Context, sid string) error
		TerminateAll(ctx context.Context, aid string) error
		TerminateWithExcept(ctx context.Context, aid, sid string) error
	}

	tokenManager interface {
		New(subject, audience string, expiration time.Duration) (string, error)
		Parse(token, audience string) (string, error)
	}
)
