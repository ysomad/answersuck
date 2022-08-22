package auth

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck-backend/internal/config"
)

type tokenService struct {
	cfg      *config.SecurityToken
	token    tokenManager
	account  accountService
	password passwordVerifier
}

type TokenServiceDeps struct {
	Config          *config.SecurityToken
	TokenManager    tokenManager
	AccountService  accountService
	PasswordMatcher passwordVerifier
}

func NewTokenService(d TokenServiceDeps) *tokenService {
	return &tokenService{
		cfg:      d.Config,
		token:    d.TokenManager,
		account:  d.AccountService,
		password: d.PasswordMatcher,
	}
}

func (s *tokenService) Create(ctx context.Context, accountId, password string) (string, error) {
	a, err := s.account.GetById(ctx, accountId)
	if err != nil {
		return "", fmt.Errorf("tokenService - Create - s.account.GetById: %w", err)
	}

	ok, err := s.password.Verify(password, a.Password)
	if err != nil {
		return "", fmt.Errorf("tokenService - Create - s.password.Verify: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("tokenService - Create - s.password.Verify: %w", ErrIncorrectAccountPassword)
	}

	t, err := s.token.Create(accountId, s.cfg.Expiration)
	if err != nil {
		return "", fmt.Errorf("tokenService - Create - s.token.Create: %w", err)
	}

	return t, nil
}

func (s *tokenService) Parse(ctx context.Context, token string) (string, error) {
	accountId, err := s.token.Parse(token)
	if err != nil {
		return "", fmt.Errorf("tokenService - Parse - s.token.Parse: %w", err)
	}

	return accountId, nil
}
