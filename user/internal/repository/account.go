package repository

import (
	"context"

	"github.com/ysomad/answersuck/user/internal/entity"

	"github.com/ysomad/answersuck/pkg/pgclient"
)

type accountRepo struct {
	client pgclient.Client
}

func NewAccountRepo() *accountRepo {
	return &accountRepo{}
}

func (r *accountRepo) Save(ctx context.Context, a *entity.Account) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepo) FindByID(ctx context.Context, accountID string) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepo) FindByEmail(ctx context.Context, email string) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepo) DeleteByID(ctx context.Context, accountID string) error {
	return nil
}
