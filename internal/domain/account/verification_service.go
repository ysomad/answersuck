package account

import (
	"context"
	"fmt"
	"time"
)

type verifService struct {
	repo  VerificationRepo
	email EmailService
}

func NewVerificationService(r VerificationRepo, e EmailService) *verifService {
	return &verifService{
		repo:  r,
		email: e,
	}
}

func (s *verifService) Request(ctx context.Context, accountId string) error {
	v, err := s.repo.Find(ctx, accountId)
	if err != nil {
		return fmt.Errorf("verifService - Request - s.repo.Find: %w", err)
	}

	if v.Verified {
		return fmt.Errorf("verifService - v.Verified: %w", ErrAlreadyVerified)
	}

	go func() {
		// TODO: handle error
		_ = s.email.SendAccountVerificationEmail(ctx, v.Email, v.Code)
	}()

	return nil
}

func (s *verifService) Verify(ctx context.Context, code string) error {
	if err := s.repo.Verify(ctx, code, time.Now()); err != nil {
		return fmt.Errorf("verifService - Verify - s.repo.Verify: %w", err)
	}

	return nil
}
