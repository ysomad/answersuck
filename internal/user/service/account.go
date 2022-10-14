package service

import (
	"context"

	"github.com/ysomad/answersuck/internal/user/domain"
	"github.com/ysomad/answersuck/internal/user/service/dto"

	"github.com/ysomad/answersuck/cryptostr"
)

type accountRepository interface {
	Save(ctx context.Context, args dto.AccountSaveArgs) (*domain.Account, error)
	FindByID(ctx context.Context, accountID string) (*domain.Account, error)
	FindByEmail(ctx context.Context, email string) (*domain.Account, error)
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

func (s *accountService) Create(ctx context.Context, email, username, plainPassword string) (*domain.Account, error) {
	// Check if username is not banned

	p, err := s.password.Encode(plainPassword)
	if err != nil {
		return nil, err
	}

	c, err := cryptostr.RandomBase64(32)
	if err != nil {
		return nil, err
	}

	a, err := s.repo.Save(ctx, dto.AccountSaveArgs{
		Email:           email,
		Username:        username,
		EncodedPassword: p,
		EmailVerifCode:  c,
	})
	if err != nil {
		return nil, err
	}

	// Send email with verification code

	return a, nil
}

func (s *accountService) GetByID(ctx context.Context, accountID string) (*domain.Account, error) {
	return nil, nil
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	return nil, nil
}

func (s *accountService) DeleteByID(ctx context.Context, accountID string) error {
	return nil
}
