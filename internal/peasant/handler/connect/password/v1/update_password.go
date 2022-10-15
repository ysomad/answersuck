package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) UpdatePassword(context.Context, *connect.Request[pb.UpdatePasswordRequest]) (*connect.Response[pb.UpdatePasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.password.v1.PasswordService.UpdatePassword is not implemented"))
}
