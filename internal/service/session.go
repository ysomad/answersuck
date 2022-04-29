package service

import (
	"context"
	"fmt"

	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
)

type sessionRepository interface {
	Create(ctx context.Context, s *domain.Session) (*domain.Session, error)
	FindById(ctx context.Context, sid string) (*domain.Session, error)
	FindAll(ctx context.Context, aid string) ([]*domain.Session, error)
	Delete(ctx context.Context, sid string) error
	DeleteAll(ctx context.Context, aid string) error
	DeleteWithExcept(ctx context.Context, aid, sid string) error
}

type session struct {
	cfg  *config.Session
	repo sessionRepository
}

func NewSession(cfg *config.Session, s sessionRepository) *session {
	return &session{
		cfg:  cfg,
		repo: s,
	}
}

func (s *session) Create(ctx context.Context, aid string, d dto.Device) (*domain.Session, error) {
	sess, err := domain.NewSession(aid, d.UserAgent, d.IP, s.cfg.Expiration)
	if err != nil {
		return nil, fmt.Errorf("session - Create - domain.NewSession: %w", err)
	}

	sess, err = s.repo.Create(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("session - Create - s.repo.Create: %w", err)
	}

	return sess, nil
}

func (s *session) GetById(ctx context.Context, sid string) (*domain.Session, error) {
	sess, err := s.repo.FindById(ctx, sid)
	if err != nil {
		return nil, fmt.Errorf("session - GetById - s.repo.FindByID: %w", err)
	}

	return sess, nil
}

func (s *session) GetAll(ctx context.Context, aid string) ([]*domain.Session, error) {
	sessions, err := s.repo.FindAll(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("session - GetAll - s.repo.FindAll: %w", err)
	}

	return sessions, nil
}

func (s *session) Terminate(ctx context.Context, sid string) error {
	if err := s.repo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("session - Terminate - s.repo.Delete: %w", err)
	}

	return nil
}

func (s *session) TerminateWithExcept(ctx context.Context, aid, sid string) error {
	if err := s.repo.DeleteWithExcept(ctx, aid, sid); err != nil {
		return fmt.Errorf("session - TerminateWithExcept - s.repo.DeleteWithExcept: %w", err)
	}

	return nil
}

func (s *session) TerminateAll(ctx context.Context, aid string) error {
	if err := s.repo.DeleteAll(ctx, aid); err != nil {
		return fmt.Errorf("session - TerminateAll - s.repo.DeleteAll: %w", err)
	}

	return nil
}
