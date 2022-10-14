package service

import (
	"context"
	"github.com/ysomad/answersuck/user/internal/entity"
	"github.com/ysomad/answersuck/user/internal/service/dto"
)

type accountRepository interface {
	Save(ctx context.Context, a *entity.Account) (*entity.Account, error)
	FindByID(ctx context.Context, accountID string) (*entity.Account, error)
	FindByEmail(ctx context.Context, email string) (*entity.Account, error)
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

func (a *accountService) Create(ctx context.Context, p dto.AccountCreateParams) (*entity.Account, error) {
	return nil, nil
}

func (a *accountService) GetByID(ctx context.Context, accountID string) (*entity.Account, error) {
	return nil, nil
}

func (a *accountService) GetByEmail(ctx context.Context, email string) (*entity.Account, error) {
	return nil, nil
}

func (a *accountService) DeleteByID(ctx context.Context, accountID string) error {
	return nil
}
