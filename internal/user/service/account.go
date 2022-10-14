package service

import (
	"context"
	"github.com/ysomad/answersuck/internal/user/service/dto"

	"github.com/ysomad/answersuck/internal/user/entity"
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

type account struct {
	repo     accountRepository
	password passwordEncodeComparer
}

func NewAccount(r accountRepository, p passwordEncodeComparer) *account {
	return &account{
		repo:     r,
		password: p,
	}
}

func (a *account) Create(ctx context.Context, p dto.AccountCreateParams) (*entity.Account, error) {
	return nil, nil
}

func (a *account) GetByID(ctx context.Context, accountID string) (*entity.Account, error) {
	return nil, nil
}

func (a *account) GetByEmail(ctx context.Context, email string) (*entity.Account, error) {
	return nil, nil
}

func (a *account) DeleteByID(ctx context.Context, accountID string) error {
	return nil
}
