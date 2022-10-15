package service

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type accountRepository interface {
	Create(ctx context.Context, args dto.AccountSaveArgs) (*domain.Account, error)
	GetByID(ctx context.Context, accountID string) (*domain.Account, error)
	DeleteByID(ctx context.Context, accountID string) error

	GetPasswordByID(ctx context.Context, accountID string) (string, error)

	UpdateEmail(ctx context.Context, accountID, newEmail string) (*domain.Account, error)
}
