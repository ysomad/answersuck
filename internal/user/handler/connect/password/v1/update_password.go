package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/password/v1"
)

func (s *server) UpdatePassword(context.Context, *connect.Request[v1.UpdatePasswordRequest]) (*connect.Response[v1.UpdatePasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.v1.PasswordService.UpdatePassword is not implemented"))
}
