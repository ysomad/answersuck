package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/cryptostr"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type emailService struct {
	accountRepo accountRepository
	verifRepo   emailVerificationRepository
	password    passwordComparer

	verifCodeLifetime time.Duration
}

func NewEmailService(ar accountRepository, vr emailVerificationRepository, p passwordComparer, lt time.Duration) (*emailService, error) {
	if ar == nil || vr == nil || p == nil || lt == 0 {
		return nil, apperror.ErrNilArgs
	}

	return &emailService{
		accountRepo:       ar,
		verifRepo:         vr,
		password:          p,
		verifCodeLifetime: lt,
	}, nil
}

func (s *emailService) Update(ctx context.Context, args dto.UpdateEmailArgs) (*domain.Account, error) {
	encoded, err := s.accountRepo.GetPasswordByID(ctx, args.AccountID)
	if err != nil {
		return nil, err
	}

	ok, err := s.password.Compare(args.PlainPassword, encoded)
	if err != nil {
		return nil, apperror.New("emailService - Update", err, domain.ErrIncorrectPassword)
	}

	if !ok {
		return nil, domain.ErrIncorrectPassword
	}

	a, err := s.accountRepo.UpdateEmail(ctx, args.AccountID, args.NewEmail)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *emailService) Verify(ctx context.Context, verifCode string) (*domain.Account, error) {
	return s.accountRepo.VerifyEmail(ctx, verifCode)
}

func (s *emailService) CreateVerification(ctx context.Context, accountID string) (domain.EmailVerification, error) {
	verifCode, err := cryptostr.RandomBase64(domain.EmailVerifCodeLen)
	if err != nil {
		return domain.EmailVerification{}, err
	}

	v := domain.NewEmailVerification(accountID, verifCode, s.verifCodeLifetime)
	if err := s.verifRepo.Save(ctx, v); err != nil {
		return domain.EmailVerification{}, err
	}

	return v, nil
}
