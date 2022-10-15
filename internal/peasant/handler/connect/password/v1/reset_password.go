package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) ResetPassword(context.Context, *connect.Request[pb.ResetPasswordRequest]) (*connect.Response[pb.ResetPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.v1.PasswordService.ResetPassword is not implemented"))
}
