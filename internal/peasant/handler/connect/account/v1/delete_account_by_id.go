package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) DeleteAccountById(context.Context, *connect.Request[pb.DeleteAccountByIdRequest]) (*connect.Response[pb.DeleteAccountByIdResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.v1.AccountService.DeleteAccount is not implemented"))
}
