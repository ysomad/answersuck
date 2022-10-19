package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type passwordService struct {
	accountRepo       accountRepository
	passwordTokenRepo passwordTokenRepository
	password          passwordEncodeComparer

	tokenLifetime time.Duration
}

func NewPasswordService(ar accountRepository, pr passwordTokenRepository, pe passwordEncodeComparer, tl time.Duration) (*passwordService, error) {
	return &passwordService{
		accountRepo:       ar,
		passwordTokenRepo: pr,
		password:          pe,
		tokenLifetime:     tl,
	}, nil
}

func (s *passwordService) Update(ctx context.Context, args dto.UpdatePasswordArgs) (*domain.Account, error) {
	oldEncodedPass, err := s.accountRepo.GetPasswordByID(ctx, args.AccountID)
	if err != nil {
		return nil, err
	}

	ok, err := s.password.Compare(args.OldPassword, oldEncodedPass)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrIncorrectPassword
	}

	newEncodedPass, err := s.password.Encode(args.NewPassword)
	if err != nil {
		return nil, err
	}

	return s.accountRepo.UpdatePassword(ctx, args.AccountID, newEncodedPass)
}

func (s *passwordService) CreateToken(ctx context.Context, emailOrUsername string) (domain.PasswordToken, error) {
	t, err := domain.GenPasswordToken()
	if err != nil {
		return domain.PasswordToken{}, err
	}

	return s.passwordTokenRepo.Create(ctx, dto.NewCreatePasswordTokenArgs(emailOrUsername, t, s.tokenLifetime))
}
