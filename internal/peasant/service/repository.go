package service

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
)

type accountRepository interface {
	Create(context.Context, dto.AccountCreateArgs, dto.EmailVerifCreateArgs) (*domain.Account, error)
	GetByID(ctx context.Context, accountID string) (*domain.Account, error)
	DeleteByID(ctx context.Context, accountID string) error

	GetPasswordByID(ctx context.Context, accountID string) (string, error)
	UpdatePassword(ctx context.Context, accountID, newPassword string) (*domain.Account, error)
	SetPassword(ctx context.Context, token, newPassword string) (*domain.Account, error)

	UpdateEmail(ctx context.Context, accountID, newEmail string) (*domain.Account, error)
	VerifyEmail(ctx context.Context, verifCode string) (*domain.Account, error)
}

type emailVerificationRepository interface {
	Save(context.Context, domain.EmailVerification) error
}

type passwordTokenRepository interface {
	Create(context.Context, dto.CreatePasswordTokenArgs) (domain.PasswordToken, error)
}