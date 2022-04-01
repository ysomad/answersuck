package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
)

type sessionService struct {
	cfg  *config.Session
	repo SessionRepo
}

func NewSessionService(cfg *config.Session, s SessionRepo) *sessionService {
	return &sessionService{
		cfg:  cfg,
		repo: s,
	}
}

func (s *sessionService) Create(ctx context.Context, aid string, d domain.Device) (*domain.Session, error) {
	sess, err := domain.NewSession(aid, d, s.cfg.Expiration)
	if err != nil {
		return nil, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	if err = s.repo.Create(ctx, sess); err != nil {
		return nil, fmt.Errorf("sessionService - Create - s.repo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetById(ctx context.Context, sid string) (*domain.Session, error) {
	sess, err := s.repo.FindById(ctx, sid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - Get - s.repo.FindByID: %w", err)
	}

	return sess, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid string) error {
	if err := s.repo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.repo.Delete: %w", err)
	}

	return nil
}
