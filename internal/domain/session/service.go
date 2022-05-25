package session

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/config"
)

type Repository interface {
	Save(ctx context.Context, s *Session) (*Session, error)
	FindById(ctx context.Context, sessionId string) (*Session, error)
	FindAll(ctx context.Context, accountId string) ([]*Session, error)
	Delete(ctx context.Context, sessionId string) error
	DeleteAll(ctx context.Context, accountId string) error
	DeleteWithExcept(ctx context.Context, accountId, sessionId string) error
}

type service struct {
	cfg  *config.Session
	repo Repository
}

func NewService(cfg *config.Session, s Repository) *service {
	return &service{
		cfg:  cfg,
		repo: s,
	}
}

func (s *service) Create(ctx context.Context, accountId string, d Device) (*Session, error) {
	sess, err := newSession(accountId, d.UserAgent, d.IP, s.cfg.Expiration)
	if err != nil {
		return nil, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	sess, err = s.repo.Save(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("sessionService - Create - s.repo.Create: %w", err)
	}

	return sess, nil
}

func (s *service) GetById(ctx context.Context, sessionId string) (*Session, error) {
	sess, err := s.repo.FindById(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetById - s.repo.FindByID: %w", err)
	}

	return sess, nil
}

func (s *service) GetAll(ctx context.Context, accountId string) ([]*Session, error) {
	sessions, err := s.repo.FindAll(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetAll - s.repo.FindAll: %w", err)
	}

	return sessions, nil
}

func (s *service) Terminate(ctx context.Context, sessionId string) error {
	if err := s.repo.Delete(ctx, sessionId); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.repo.Delete: %w", err)
	}

	return nil
}

func (s *service) TerminateWithExcept(ctx context.Context, accountId, sessionId string) error {
	if err := s.repo.DeleteWithExcept(ctx, accountId, sessionId); err != nil {
		return fmt.Errorf("sessionService - TerminateWithExcept - s.repo.DeleteWithExcept: %w", err)
	}

	return nil
}

func (s *service) TerminateAll(ctx context.Context, accountId string) error {
	if err := s.repo.DeleteAll(ctx, accountId); err != nil {
		return fmt.Errorf("sessionService - TerminateAll - s.repo.DeleteAll: %w", err)
	}

	return nil
}
