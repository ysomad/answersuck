package auth

import (
	"context"
	"time"

	"github.com/ysomad/answersuck-backend/internal/domain/account"
	"github.com/ysomad/answersuck-backend/internal/domain/session"
)

type (
	accountService interface {
		GetById(ctx context.Context, accountId string) (account.Account, error)
		GetByEmail(ctx context.Context, email string) (account.Account, error)
		GetByNickname(ctx context.Context, nickname string) (account.Account, error)
	}

	sessionService interface {
		Create(ctx context.Context, accountId string, d session.Device) (*session.Session, error)
		Terminate(ctx context.Context, sessionId string) error
	}

	tokenManager interface {
		Create(subject string, expiration time.Duration) (string, error)
		Parse(token string) (string, error)
	}

	passwordVerifier interface {
		Verify(plain, hash string) (bool, error)
	}
)
