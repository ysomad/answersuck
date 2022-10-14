package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/account/v1"
)

func (s *server) GetAccountByEmail(context.Context, *connect.Request[v1.GetAccountByEmailRequest]) (*connect.Response[v1.GetAccountByEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.v1.AccountService.GetAccountByEmail is not implemented"))
}
