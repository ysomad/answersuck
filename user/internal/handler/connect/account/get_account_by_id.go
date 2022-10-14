package account

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/user/internal/gen/proto/account"
)

func (s *server) GetAccountById(ctx context.Context, r *connect.Request[account.GetAccountByIdRequest]) (*connect.Response[account.GetAccountByIdResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.AccountService.GetAccountById is not implemented"))
}
