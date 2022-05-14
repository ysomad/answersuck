package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) (*domain.Session, error)
	FindById(ctx context.Context, sessionId string) (*domain.Session, error)
	FindAll(ctx context.Context, accountId string) ([]*domain.Session, error)
	Delete(ctx context.Context, sessionId string) error
	DeleteAll(ctx context.Context, accountId string) error
	DeleteWithExcept(ctx context.Context, accountId, sessionId string) error
}

type sessionService struct {
	cfg  *config.Session
	repo SessionRepository
}

func NewSessionService(cfg *config.Session, s SessionRepository) *sessionService {
	return &sessionService{
		cfg:  cfg,
		repo: s,
	}
}

func (s *sessionService) Create(ctx context.Context, accountId string, d domain.Device) (*domain.Session, error) {
	sess, err := domain.NewSession(accountId, d.UserAgent, d.IP, s.cfg.Expiration)
	if err != nil {
		return nil, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	sess, err = s.repo.Create(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("sessionService - Create - s.repo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetById(ctx context.Context, sessionId string) (*domain.Session, error) {
	sess, err := s.repo.FindById(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetById - s.repo.FindByID: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetAll(ctx context.Context, accountId string) ([]*domain.Session, error) {
	sessions, err := s.repo.FindAll(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetAll - s.repo.FindAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sessionId string) error {
	if err := s.repo.Delete(ctx, sessionId); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.repo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateWithExcept(ctx context.Context, accountId, sessionId string) error {
	if err := s.repo.DeleteWithExcept(ctx, accountId, sessionId); err != nil {
		return fmt.Errorf("sessionService - TerminateWithExcept - s.repo.DeleteWithExcept: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, accountId string) error {
	if err := s.repo.DeleteAll(ctx, accountId); err != nil {
		return fmt.Errorf("sessionService - TerminateAll - s.repo.DeleteAll: %w", err)
	}

	return nil
}
