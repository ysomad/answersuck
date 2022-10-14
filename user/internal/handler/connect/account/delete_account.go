package account

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ysomad/answersuck/user/internal/gen/proto/account"
)

func (s *server) DeleteAccount(ctx context.Context, r *connect.Request[account.DeleteAccountRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.account.AccountService.DeleteAccount is not implemented"))
}
