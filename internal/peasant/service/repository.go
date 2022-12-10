package service

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type accountRepository interface {
	Create(context.Context, dto.AccountCreateArgs) (*domain.Account, error)
	GetByID(ctx context.Context, accountID string) (*domain.Account, error)
	GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (*domain.Account, error)
	DeleteByID(ctx context.Context, accountID string) error

	GetPasswordByID(ctx context.Context, accountID string) (string, error)
	UpdatePassword(ctx context.Context, accountID, newPassword string) (*domain.Account, error)

	UpdateEmail(ctx context.Context, accountID, newEmail string) (*domain.Account, error)
	VerifyEmail(ctx context.Context, accountID string) (*domain.Account, error)
}
