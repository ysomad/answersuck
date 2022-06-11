package auth

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"
)

type tokenService struct {
	cfg     *config.AccessToken
	token   TokenManager
	account AccountService
}

func NewTokenService(c *config.AccessToken, t TokenManager, a AccountService) *tokenService {
	return &tokenService{
		cfg:     c,
		token:   t,
		account: a,
	}
}

func (s *tokenService) Create(ctx context.Context, dto TokenCreateDTO) (string, error) {
	a, err := s.account.GetById(ctx, dto.AccountId)
	if err != nil {
		return "", fmt.Errorf("tokenService - Create - s.account.GetById: %w", err)
	}

	if err = a.ComparePasswords(dto.Password); err != nil {
		return "", fmt.Errorf("tokenService - Create - a.ComparePasswords: %w", err)
	}

	t, err := s.token.Create(dto.AccountId, dto.Audience, s.cfg.Expiration)
	if err != nil {
		return "", fmt.Errorf("tokenService - Create - s.token.New: %w", err)
	}

	return t, nil
}

func (s *tokenService) Parse(ctx context.Context, token, audience string) (string, error) {
	accountId, err := s.token.Parse(token, audience)
	if err != nil {
		return "", fmt.Errorf("tokenService - Parse - s.token.Parse: %w", err)
	}

	return accountId, nil
}
