package service

import (
	"context"
	"fmt"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/domain"
)

type sessionService struct {
	cfg  *app.Config
	repo SessionRepo
}

func NewSessionService(cfg *app.Config, s SessionRepo) *sessionService {
	return &sessionService{
		cfg:  cfg,
		repo: s,
	}
}

func (s *sessionService) Create(ctx context.Context, aid string, d domain.Device) (domain.Session, error) {
	sess, err := domain.NewSession(aid, d, s.cfg.SessionTTL)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - domain.NewSession: %w", err)
	}

	if err = s.repo.Create(ctx, sess); err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Create - s.repo.Create: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetByID(ctx context.Context, sid string) (domain.Session, error) {
	sess, err := s.repo.FindByID(ctx, sid)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessionService - Get - s.repo.FindByID: %w", err)
	}

	return sess, nil
}

func (s *sessionService) GetAll(ctx context.Context, aid string) ([]domain.Session, error) {
	sessions, err := s.repo.FindAll(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("sessionService - GetAll - s.repo.FindAll: %w", err)
	}

	return sessions, nil
}

func (s *sessionService) Terminate(ctx context.Context, sid, currSid string) error {
	if sid == currSid {
		return fmt.Errorf("sessionService - Terminate: %w", domain.ErrSessionNotTerminated)
	}

	if err := s.repo.Delete(ctx, sid); err != nil {
		return fmt.Errorf("sessionService - Terminate - s.repo.Delete: %w", err)
	}

	return nil
}

func (s *sessionService) TerminateAll(ctx context.Context, aid, sid string) error {
	if err := s.repo.DeleteAll(ctx, aid, sid); err != nil {
		return fmt.Errorf("sessionService - TerminateAll - s.repo.DeleteAll: %w", err)
	}

	return nil
}
