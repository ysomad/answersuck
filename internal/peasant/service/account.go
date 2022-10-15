package service

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"

	"github.com/ysomad/answersuck/cryptostr"
)

type accountRepository interface {
	Save(ctx context.Context, args dto.AccountSaveArgs) (*domain.Account, error)
	FindByID(ctx context.Context, accountID string) (*domain.Account, error)
	DeleteByID(ctx context.Context, accountID string) error
}

type passwordEncodeComparer interface {
	Encode(plain string) (string, error)
	Compare(plain, encoded string) (bool, error)
}

type accountService struct {
	repo     accountRepository
	password passwordEncodeComparer
}

func NewAccountService(r accountRepository, p passwordEncodeComparer) *accountService {
	return &accountService{
		repo:     r,
		password: p,
	}
}

func (s *accountService) Create(ctx context.Context, args dto.AccountCreateArgs) (*domain.Account, error) {
	// Check if username is not banned

	encodedPass, err := s.password.Encode(args.PlainPassword)
	if err != nil {
		return nil, err
	}

	verifCode, err := cryptostr.RandomBase64(32)
	if err != nil {
		return nil, err
	}

	a, err := s.repo.Save(ctx, dto.AccountSaveArgs{
		Email:           args.Email,
		Username:        args.Username,
		EncodedPassword: encodedPass,
		EmailVerifCode:  verifCode,
	})
	if err != nil {
		return nil, err
	}

	// Send email with verification code

	return a, nil
}

func (s *accountService) GetByID(ctx context.Context, accountID string) (*domain.Account, error) {
	return s.repo.FindByID(ctx, accountID)
}

func (s *accountService) DeleteByID(ctx context.Context, accountID string) error {
	return s.repo.DeleteByID(ctx, accountID)
}
