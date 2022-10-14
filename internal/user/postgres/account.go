package postgres

import (
	"context"

	"github.com/ysomad/answersuck/internal/user/entity"

	"github.com/ysomad/answersuck/pgclient"
)

type accountRepository struct {
	client *pgclient.Client
}

func NewAccountRepository(c *pgclient.Client) *accountRepository {
	return &accountRepository{
		client: c,
	}
}

func (r *accountRepository) Save(ctx context.Context, a *entity.Account) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepository) FindByID(ctx context.Context, accountID string) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*entity.Account, error) {
	return nil, nil
}

func (r *accountRepository) DeleteByID(ctx context.Context, accountID string) error {
	return nil
}
