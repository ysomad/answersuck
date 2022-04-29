package service

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
)

type accountService interface {
	Create(ctx context.Context, a *domain.Account) (*domain.Account, error)
	GetById(ctx context.Context, aid string) (*domain.Account, error)
	GetByEmail(ctx context.Context, email string) (*domain.Account, error)
	GetByUsername(ctx context.Context, username string) (*domain.Account, error)
	Delete(ctx context.Context, aid, sid string) error
	RequestVerification(ctx context.Context, aid string) error
	Verify(ctx context.Context, code string, verified bool) error
	RequestPasswordReset(ctx context.Context, login string) error
	PasswordReset(ctx context.Context, token, password string) error
}

type auth struct {
	cfg     *config.Aggregate
	token   tokenManager
	account accountService
	session sessionService
}

func NewAuth(cfg *config.Aggregate, t tokenManager, a accountService, s sessionService) *auth {
	return &auth{
		cfg:     cfg,
		token:   t,
		account: a,
		session: s,
	}
}

func (s *auth) Login(ctx context.Context, login, password string, d dto.Device) (*domain.Session, error) {
	a := &domain.Account{}

	_, err := mail.ParseAddress(login)
	if err != nil {
		// login is not email
		a, err = s.account.GetByUsername(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("auth - Login - s.account.GetByUsername: %w", err)
		}
	} else {
		// login is email
		a, err = s.account.GetByEmail(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("auth - Login - s.account.GetByEmail: %w", err)
		}
	}

	a.Password = password

	if err = a.CompareHashAndPassword(); err != nil {
		return nil, fmt.Errorf("auth - Login - a.CompareHashAndPassword: %w", err)
	}

	sess, err := s.session.Create(ctx, a.Id, d)
	if err != nil {
		return nil, fmt.Errorf("auth - Login - s.session.Create: %w", err)
	}

	return sess, nil
}

func (s *auth) Logout(ctx context.Context, sid string) error {
	if err := s.session.Terminate(ctx, sid); err != nil {
		return fmt.Errorf("auth - Logout - s.session.Terminate: %w", err)
	}

	return nil
}

func (s *auth) NewAccessToken(ctx context.Context, aid, password, audience string) (string, error) {
	a, err := s.account.GetById(ctx, aid)
	if err != nil {
		return "", fmt.Errorf("auth - NewAccessToken - s.account.GetByID: %w", err)
	}

	a.Password = password

	if err = a.CompareHashAndPassword(); err != nil {
		return "", fmt.Errorf("auth - NewAccessToken - a.CompareHashAndPassword: %w", err)
	}

	t, err := s.token.New(aid, audience, s.cfg.AccessToken.Expiration)
	if err != nil {
		return "", fmt.Errorf("auth - NewAccessToken - s.token.New: %w", err)
	}

	return t, nil
}

func (s *auth) ParseAccessToken(ctx context.Context, token, audience string) (string, error) {
	aid, err := s.token.Parse(token, audience)
	if err != nil {
		return "", fmt.Errorf("auth - ParseAccessToken - s.token.Parse: %w", err)
	}

	return aid, nil
}
