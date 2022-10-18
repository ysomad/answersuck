package service

import (
	"context"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type passwordService struct {
	accountRepo accountRepository
	password    passwordEncodeComparer
}

func NewPasswordService(ar accountRepository, pe passwordEncodeComparer) (*passwordService, error) {
	if ar == nil || pe == nil {
		return nil, apperror.ErrNilArgs
	}

	return &passwordService{
		accountRepo: ar,
		password:    pe,
	}, nil
}

func (s *passwordService) Update(ctx context.Context, args dto.UpdatePasswordArgs) (*domain.Account, error) {
	oldEncodedPassword, err := s.accountRepo.GetPasswordByID(ctx, args.AccountID)
	if err != nil {
		return nil, err
	}

	match, err := s.password.Compare(args.OldPassword, oldEncodedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, domain.ErrIncorrectPassword
	}

	newEncodedPassword, err := s.password.Encode(args.NewPassword)
	if err != nil {
		return nil, err
	}

	return s.accountRepo.UpdatePassword(ctx, args.AccountID, newEncodedPassword)
}
