package account

import (
	"context"
	"time"
)

type repository interface {
	Save(ctx context.Context, a Account, code string) (Account, error)
	FindById(ctx context.Context, accountId string) (Account, error)
	FindByEmail(ctx context.Context, email string) (Account, error)
	FindByNickname(ctx context.Context, nickname string) (Account, error)
	SetArchived(ctx context.Context, accountId string, archived bool, updatedAt time.Time) error

	Verify(ctx context.Context, code string, updatedAt time.Time) error
	FindVerification(ctx context.Context, nickname string) (Verification, error)

	SavePasswordToken(ctx context.Context, dto SavePasswordTokenDTO) (email string, err error)
	FindPasswordToken(ctx context.Context, token string) (PasswordToken, error)
	SetPassword(ctx context.Context, dto SetPasswordDTO) error
}

type sessionService interface {
	TerminateAll(ctx context.Context, accountId string) error
}

type emailService interface {
	SendAccountVerificationEmail(ctx context.Context, to, code string) error
	SendPasswordResetEmail(ctx context.Context, to, token string) error
}

type blockList interface {
	Find(nickname string) (found bool)
}

type passwordHasher interface {
	Hash(plain string) (hash string, err error)
}
