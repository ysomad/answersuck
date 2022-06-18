package account

import (
	"context"
	"time"
)

type AccountRepo interface {
	Save(ctx context.Context, a *Account, code string) (string, error)
	FindEmailByNickname(ctx context.Context, nickname string) (string, error)
	FindById(ctx context.Context, accountId string) (*Account, error)
	FindByEmail(ctx context.Context, email string) (*Account, error)
	FindByNickname(ctx context.Context, nickname string) (*Account, error)
	Archive(ctx context.Context, accountId string, updatedAt time.Time) error

	SavePasswordToken(ctx context.Context, email, token string) error
	FindPasswordToken(ctx context.Context, token string) (PasswordToken, error)
	SetPassword(ctx context.Context, dto SetPasswordDTO) error
}

type VerificationRepo interface {
	Verify(ctx context.Context, code string, updatedAt time.Time) error
	Find(ctx context.Context, nickname string) (Verification, error)
}
