package account

import "context"

type SessionService interface {
	TerminateAll(ctx context.Context, accountId string) error
}

type EmailService interface {
	SendAccountVerificationEmail(ctx context.Context, to, code string) error
	SendPasswordResetEmail(ctx context.Context, to, token string) error
}

type BlockList interface {
	Find(s string) bool
}
