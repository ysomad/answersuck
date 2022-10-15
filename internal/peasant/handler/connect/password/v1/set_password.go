package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) SetPassword(context.Context, *connect.Request[pb.SetPasswordRequest]) (*connect.Response[pb.SetPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.v1.PasswordService.SetPassword is not implemented"))
}
