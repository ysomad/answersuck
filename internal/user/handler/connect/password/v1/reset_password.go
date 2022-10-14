package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/password/v1"
)

func (s *server) ResetPassword(context.Context, *connect.Request[v1.ResetPasswordRequest]) (*connect.Response[v1.ResetPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.v1.PasswordService.ResetPassword is not implemented"))
}
