package service

import (
	"context"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type passwordVerifier interface {
	Compare(plain, encoded string) (bool, error)
}

type emailService struct {
	accountRepo accountRepository
	password    passwordVerifier
}

func NewEmailService(r accountRepository, p passwordVerifier) (*emailService, error) {
	if r == nil || p == nil {
		return nil, apperror.ErrNilArgs
	}

	return &emailService{
		accountRepo: r,
		password:    p,
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

	// TODO: send email to verify email

	return a, nil
}

func (s *emailService) Verify(ctx context.Context, verifCode string) (*domain.Account, error) {
	return s.accountRepo.VerifyEmail(ctx, verifCode)
}

// TODO: REFACTOR NAMING??!!
func (s *emailService) SendVerification(ctx context.Context, accountID string) error {
	// TODO: implement emailService.SendVerification
	return nil
}
