package service

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
)

type AccountService interface {
	GetById(ctx context.Context, accountId string) (*domain.Account, error)
	GetByEmail(ctx context.Context, email string) (*domain.Account, error)
	GetByUsername(ctx context.Context, username string) (*domain.Account, error)
}

type authService struct {
	cfg     *config.Aggregate
	token   TokenManager
	account AccountService
	session SessionService
}

func NewAuthService(cfg *config.Aggregate, t TokenManager, a AccountService, s SessionService) *authService {
	return &authService{
		cfg:     cfg,
		token:   t,
		account: a,
		session: s,
	}
}

func (s *authService) Login(ctx context.Context, login, password string, d domain.Device) (*domain.Session, error) {
	var a *domain.Account

	_, err := mail.ParseAddress(login)
	if err != nil {
		a, err = s.account.GetByUsername(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("authService - Login - s.account.GetByUsername: %w", err)
		}
	} else {
		a, err = s.account.GetByEmail(ctx, login)
		if err != nil {
			return nil, fmt.Errorf("authService - Login - s.account.GetByEmail: %w", err)
		}
	}

	a.Password = password

	if err = a.CompareHashAndPassword(); err != nil {
		return nil, fmt.Errorf("authService - Login - a.CompareHashAndPassword: %w", err)
	}

	sess, err := s.session.Create(ctx, a.Id, d)
	if err != nil {
		return nil, fmt.Errorf("authService - Login - s.session.Create: %w", err)
	}

	return sess, nil
}

func (s *authService) Logout(ctx context.Context, sessionId string) error {
	if err := s.session.Terminate(ctx, sessionId); err != nil {
		return fmt.Errorf("authService - Logout - s.session.Terminate: %w", err)
	}

	return nil
}

func (s *authService) NewSecurityToken(ctx context.Context, accountId, password, audience string) (string, error) {
	a, err := s.account.GetById(ctx, accountId)
	if err != nil {
		return "", fmt.Errorf("authService - NewSecurityToken - s.account.GetByID: %w", err)
	}

	a.Password = password

	if err = a.CompareHashAndPassword(); err != nil {
		return "", fmt.Errorf("authService - NewSecurityToken - a.CompareHashAndPassword: %w", err)
	}

	t, err := s.token.New(accountId, audience, s.cfg.AccessToken.Expiration)
	if err != nil {
		return "", fmt.Errorf("authService - NewSecurityToken - s.token.New: %w", err)
	}

	return t, nil
}

func (s *authService) ParseSecurityToken(ctx context.Context, token, audience string) (string, error) {
	accountId, err := s.token.Parse(token, audience)
	if err != nil {
		return "", fmt.Errorf("authService - ParseSecurityToken - s.token.Parse: %w", err)
	}

	return accountId, nil
}
