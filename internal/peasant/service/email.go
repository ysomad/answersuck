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
	repo     accountRepository
	password passwordVerifier
}

func NewEmailService(r accountRepository, p passwordVerifier) (*emailService, error) {
	if r == nil || p == nil {
		return nil, apperror.ErrNilArgs
	}

	return &emailService{
		repo:     r,
		password: p,
	}, nil
}

func (s *emailService) Update(ctx context.Context, args dto.UpdateEmailArgs) (*domain.Account, error) {
	encoded, err := s.repo.GetPasswordByID(ctx, args.AccountID)
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

	return s.repo.UpdateEmail(ctx, args.AccountID, args.NewEmail)
}

// TODO: REFACTOR NAMING??!!
func (s *emailService) SendVerification(ctx context.Context, accountID string) error {
	return nil
}

func (s *emailService) Verify(ctx context.Context, code string) (*domain.Account, error) {

	return nil, nil
}
