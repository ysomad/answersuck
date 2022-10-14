package account

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/user/internal/gen/proto/account"
)

func (s *server) GetAccountByEmail(ctx context.Context, r *connect.Request[account.GetAccountByEmailRequest]) (*connect.Response[account.GetAccountByEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.AccountService.GetAccountByEmail is not implemented"))
}
