package accountv1

import (
	"context"
	pb "github.com/ysomad/answersuck-backend/api/proto/account/v1"
	"github.com/ysomad/answersuck-backend/internal/domain/account"
)

var _ pb.AccountService = &server{}

type accountService interface {
	Create(ctx context.Context, email, nickname, password string) (account.Account, error)
	Archive(ctx context.Context, accountID string) error
	RequestEmailVerification(ctx context.Context, accountID string) error
	VerifyEmail(ctx context.Context, code string) (account.Account, error)
	ResetPassword(ctx context.Context, login string) error
	SetPassword(ctx context.Context, token, password string) error
	UpdatePassword(ctx context.Context, args account.UpdatePasswordParams) error
}

type server struct {
	account accountService
}

func NewServer(as accountService) *server {
	return &server{
		account: as,
	}
}
