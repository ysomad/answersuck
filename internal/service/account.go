package service

import (
	"context"
	"fmt"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/domain"

	"github.com/quizlyfun/quizly-backend/pkg/auth"
)

type accountService struct {
	cfg     *app.Config
	repo    AccountRepo
	token   auth.TokenManager
	session Session
	email   Email
}

func NewAccountService(cfg *app.Config, r AccountRepo, s Session, t auth.TokenManager, e Email) *accountService {
	return &accountService{
		cfg:     cfg,
		repo:    r,
		token:   t,
		session: s,
		email:   e,
	}
}

func (s *accountService) Create(ctx context.Context, acc domain.Account) (domain.Account, error) {
	if err := acc.GeneratePasswordHash(); err != nil {
		return domain.Account{}, fmt.Errorf("accountService - Create - acc.GeneratePasswordHash: %w", err)
	}

	// TODO: implement email verification
	acc.Verified = true

	a, err := s.repo.Create(ctx, acc)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - Create - s.repo.Create: %w", err)
	}

	if err = s.email.SendAccountVerificationEmail(ctx, a.Email); err != nil {
		return domain.Account{}, fmt.Errorf("accountService - Create - s.email.SendVerification: %w", err)
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

func (s *accountService) GetByUsername(ctx context.Context, uname string) (domain.Account, error) {
	acc, err := s.repo.FindByUsername(ctx, uname)
	if err != nil {
		return domain.Account{}, fmt.Errorf("accountService - GetByUsername - s.repo.FindByUsername: %w", err)
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

func (s *accountService) Verify(ctx context.Context, aid, code string) error {
	panic("implement")

	return nil
}
