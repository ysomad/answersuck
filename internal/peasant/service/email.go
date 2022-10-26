package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type emailService struct {
	accountRepo accountRepository
	password    passwordComparer

	verifCodeLifetime time.Duration
}

func NewEmailService(ar accountRepository, p passwordComparer, lt time.Duration) (*emailService, error) {
	return &emailService{
		accountRepo:       ar,
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

// TODO: reimplement
func (s *emailService) Verify(ctx context.Context, token string) (*domain.Account, error) {
	// verify token and get account id from it
	accountID := ""
	return s.accountRepo.VerifyEmail(ctx, accountID)
}

// NotifyWithToken creates new email verification token and notifies user with it in url, user must visit
// the url to verify email.
func (s *emailService) NotifyWithToken(ctx context.Context, accountID string) (domain.EmailVerifToken, error) {
	t := domain.NewEmailVerifToken(accountID)

	// TODO: create new email verif token
	// TODO: send email verif token to user email

	return t, nil
}
