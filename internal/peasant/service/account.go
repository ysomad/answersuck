package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type accountService struct {
	repo     accountRepository
	password passwordEncodeComparer

	verifTokenLtime time.Duration
}

func NewAccountService(r accountRepository, ec passwordEncodeComparer, tl time.Duration) (*accountService, error) {
	return &accountService{
		repo:            r,
		password:        ec,
		verifTokenLtime: tl,
	}, nil
}

func (s *accountService) Create(ctx context.Context, args dto.AccountCreateArgs) (a *domain.Account, err error) {
	// TODO: Check if password is not banned

	// TODO: Check if username is not banned

	// TODO: Check if email is real or not banned

	args.Password, err = s.password.Encode(args.Password)
	if err != nil {
		return nil, err
	}

	a, err = s.repo.Create(ctx, args)
	if err != nil {
		return nil, err
	}

	// TODO: create new email verification token
	_ = domain.NewEmailVerifToken(a.ID)

	// TODO: Send email with verification token

	return a, nil
}

func (s *accountService) GetByID(ctx context.Context, accountID string) (*domain.Account, error) {
	return s.repo.GetByID(ctx, accountID)
}

func (s *accountService) DeleteByID(ctx context.Context, accountID, password string) error {
	encodedPass, err := s.repo.GetPasswordByID(ctx, accountID)
	if err != nil {
		return err
	}

	ok, err := s.password.Compare(password, encodedPass)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrIncorrectPassword
	}

	// TODO: log out all sessions
	return s.repo.DeleteByID(ctx, accountID)
}
