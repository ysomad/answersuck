package service

import (
	"context"

	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/internal/dto"
)

type (
	Account interface {
		// Create creates new account
		Create(ctx context.Context, acc *domain.Account) (*domain.Account, error)

		// GetById account
		GetById(ctx context.Context, aid string) (*domain.Account, error)

		// GetByEmail account
		GetByEmail(ctx context.Context, email string) (*domain.Account, error)

		// GetByUsername account
		GetByUsername(ctx context.Context, username string) (*domain.Account, error)

		// Delete sets account IsArchive state to true
		Delete(ctx context.Context, aid, sid string) error

		// Verify verifies account using provided token
		Verify(ctx context.Context, aid, code string, verified bool) error
	}

	AccountRepo interface {
		// Create account
		Create(ctx context.Context, a *domain.Account) (*domain.Account, error)

		// FindById account in DB
		FindById(ctx context.Context, aid string) (*domain.Account, error)

		// FindByEmail account in DB
		FindByEmail(ctx context.Context, email string) (*domain.Account, error)

		// FindByUsername account in DB
		FindByUsername(ctx context.Context, username string) (*domain.Account, error)

		// Archive sets entity.Account.IsArchive state to provided value
		Archive(ctx context.Context, aid string, archive bool) error

		// Verify sets verified to account with code in account_verification entity
		Verify(ctx context.Context, a dto.AccountVerification) error
	}

	Auth interface {
		// Login creates new session using provided account email or username as login and password
		Login(ctx context.Context, login, password string, d domain.Device) (*domain.Session, error)

		// Logout logs out session by id
		Logout(ctx context.Context, sid string) error

		// NewAccessToken creates new token with subject (accountID) and audience claims
		NewAccessToken(ctx context.Context, aid, password, audience string) (string, error)

		// ParseAccessToken parses token and returns subject from claims
		ParseAccessToken(ctx context.Context, token, audience string) (string, error)
	}

	Email interface {
		// SendAccountVerificationEmail sends email with verification link to given email as to.
		SendAccountVerificationEmail(ctx context.Context, to, username, code string) error
	}

	Session interface {
		// Create new session for account with id and device of given provider
		Create(ctx context.Context, aid string, d domain.Device) (*domain.Session, error)

		// GetById session
		GetById(ctx context.Context, sid string) (*domain.Session, error)

		// Terminate session by id
		Terminate(ctx context.Context, sid string) error
	}

	SessionRepo interface {
		// Create new session in DB
		Create(ctx context.Context, s *domain.Session) error

		// FindById session
		FindById(ctx context.Context, sid string) (*domain.Session, error)

		// Delete session by id
		Delete(ctx context.Context, sid string) error
	}
)
