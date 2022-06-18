package v1

import (
	"context"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/auth"
	"github.com/answersuck/vault/internal/domain/media"
	"github.com/answersuck/vault/internal/domain/session"
)

type (
	AccountService interface {
		Create(ctx context.Context, r account.CreateReq) (*account.Account, error)
		Delete(ctx context.Context, accountId string) error

		ResetPassword(ctx context.Context, login string) error
		SetPassword(ctx context.Context, token, password string) error
	}

	VerificationService interface {
		Request(ctx context.Context, accountId string) error
		Verify(ctx context.Context, code string) error
	}
)

type SessionService interface {
	GetByIdWithVerified(ctx context.Context, sessionId string) (*session.WithAccountDetails, error)
	GetById(ctx context.Context, sessionId string) (*session.Session, error)
	GetAll(ctx context.Context, accountId string) ([]*session.Session, error)
	Terminate(ctx context.Context, sessionId string) error
	TerminateWithExcept(ctx context.Context, accountId, sessionId string) error
}

type (
	LoginService interface {
		Login(ctx context.Context, login, password string, d session.Device) (*session.Session, error)
	}

	TokenService interface {
		Create(ctx context.Context, dto auth.TokenCreateDTO) (string, error)
		Parse(ctx context.Context, token, audience string) (string, error)
	}
)

type MediaService interface {
	UploadAndSave(ctx context.Context, dto *media.UploadDTO) (media.Media, error)
}
