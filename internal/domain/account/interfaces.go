package account

import (
	"context"
	"time"
)

type AccountRepo interface {
	Save(ctx context.Context, a Account, code string) (Account, error)
	FindById(ctx context.Context, accountId string) (Account, error)
	FindByEmail(ctx context.Context, email string) (Account, error)
	FindByNickname(ctx context.Context, nickname string) (Account, error)
	Archive(ctx context.Context, accountId string, updatedAt time.Time) error

	Verify(ctx context.Context, code string, updatedAt time.Time) error
	FindVerification(ctx context.Context, nickname string) (Verification, error)

	// SavePasswordToken saves password token for account with login, returns email
	SavePasswordToken(ctx context.Context, dto SavePasswordTokenDTO) (email string, err error)
	FindPasswordToken(ctx context.Context, token string) (PasswordToken, error)
	SetPassword(ctx context.Context, dto SetPasswordDTO) error
}

type SessionService interface {
	TerminateAll(ctx context.Context, accountId string) error
}

type EmailService interface {
	SendAccountVerificationEmail(ctx context.Context, to, code string) error
	SendPasswordResetEmail(ctx context.Context, to, token string) error
}

type BlockList interface {
	Find(nickname string) bool
}

type Password interface {
	Hash(plain string) (string, error)
}
