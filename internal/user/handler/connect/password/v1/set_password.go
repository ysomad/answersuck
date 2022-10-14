package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	v1 "github.com/ysomad/answersuck/rpc/user/password/v1"
)

func (s *server) SetPassword(context.Context, *connect.Request[v1.SetPasswordRequest]) (*connect.Response[v1.SetPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.v1.PasswordService.SetPassword is not implemented"))
}
