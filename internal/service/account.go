package service

import (
	"context"
	"fmt"

	"github.com/Quizish/quizish-backend/internal/app"
	"github.com/Quizish/quizish-backend/internal/domain"
)

type accountService struct {
	cfg     *app.Config
	repo    AccountRepo
	session Session
}

func NewAccountService(cfg *app.Config, r AccountRepo, s Session) *accountService {
	return &accountService{
		cfg:     cfg,
		repo:    r,
		session: s,
	}
}

func (s *accountService) Create(ctx context.Context, acc domain.Account) (domain.Account, error) {
	if err := acc.GeneratePasswordHash(); err != nil {
		return domain.Account{}, fmt.Errorf("accountService - Create - acc.GeneratePasswordHash: %w", err)
	}

	a, err := s.repo.Create(ctx, acc)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - Create - s.repo.Create: %w", err)
	}

	return a, nil
}

func (s *accountService) GetByID(ctx context.Context, aid string) (domain.Account, error) {
	acc, err := s.repo.FindByID(ctx, aid)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByID - s.repo.FindByID: %w", err)
	}

	return acc, nil
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (domain.Account, error) {
	acc, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByEmail - s.repo.FindByEmail: %w", err)
	}

	return acc, nil
}

func (s *accountService) Delete(ctx context.Context, aid, sid string) error {
	if err := s.repo.Archive(ctx, aid, true); err != nil {
		return fmt.Errorf("accountService - Archive - s.repo.Archive: %w", err)
	}

	if err := s.session.TerminateAll(ctx, aid, sid); err != nil {
		return fmt.Errorf("accountService - Archive - s.session.TerminateAll: %w", err)
	}

	return nil
}

func (s *accountService) Verify(ctx context.Context, code string) error {
	panic("implement")

	return nil
}
