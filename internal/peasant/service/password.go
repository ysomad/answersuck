package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type passwordService struct {
	accountRepo accountRepository
	password    passwordEncodeComparer

	tokenLtime time.Duration
}

func NewPasswordService(ar accountRepository, ec passwordEncodeComparer, tl time.Duration) (*passwordService, error) {
	return &passwordService{
		accountRepo: ar,
		password:    ec,
		tokenLtime:  tl,
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

// NotifyWithToken finds account with email or username and creates password token with found account id,
// then notifies user with token in url, user must to visit the url to set new password.
func (s *passwordService) NotifyWithToken(ctx context.Context, emailOrUsername string) (domain.PasswordToken, error) {
	// TODO: get account by email or username
	// TODO: implement password service create token

	t := domain.NewPasswordToken()

	return t, nil
}

func (s *passwordService) Set(ctx context.Context, token, newPassword string) (*domain.Account, error) {
	// TODO: verify token and get accountID from it
	accountID := ""

	newEncodedPass, err := s.password.Encode(newPassword)
	if err != nil {
		return nil, err
	}

	return s.accountRepo.UpdatePassword(ctx, accountID, newEncodedPass)
}
