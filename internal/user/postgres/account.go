package postgres

import (
	"context"

	"github.com/ysomad/answersuck/internal/user/entity"
)

type accountRepository struct{}

func NewAccountRepository() *accountRepository {
	return &accountRepository{}
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
