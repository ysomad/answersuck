package auth

import (
	"context"
	"time"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"
)

type (
	AccountService interface {
		GetById(ctx context.Context, accountId string) (account.Account, error)
		GetByEmail(ctx context.Context, email string) (account.Account, error)
		GetByNickname(ctx context.Context, nickname string) (account.Account, error)
	}

	SessionService interface {
		Create(ctx context.Context, accountId string, d session.Device) (*session.Session, error)
		Terminate(ctx context.Context, sessionId string) error
	}

	TokenManager interface {
		Create(subject string, expiration time.Duration) (string, error)
		Parse(token string) (string, error)
	}

	PasswordVerifier interface {
		Verify(plain, hash string) (bool, error)
	}
)
