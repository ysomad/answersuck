package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/account/v1"
)

func (s *server) DeleteAccount(context.Context, *connect.Request[v1.DeleteAccountRequest]) (*connect.Response[v1.DeleteAccountResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.v1.AccountService.DeleteAccount is not implemented"))
}
