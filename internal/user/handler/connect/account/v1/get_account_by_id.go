package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/account/v1"
)

func (s *server) GetAccountById(ctx context.Context, r *connect.Request[v1.GetAccountByIdRequest]) (*connect.Response[v1.GetAccountByIdResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.v1.AccountService.GetAccountById is not implemented"))
}
