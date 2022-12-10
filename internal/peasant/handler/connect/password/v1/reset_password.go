package v1

import (
	"context"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) ResetPassword(ctx context.Context, r *connect.Request[pb.ResetPasswordRequest]) (*connect.Response[pb.ResetPasswordResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	t, err := s.passwordService.CreateSetterToken(ctx, r.Msg.GetEmailOrUsername())
	if err != nil {
		s.log.Error(err.Error())

		// TODO: ResetPassword handle specific errors

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.ResetPasswordResponse{
		Token: t.String(),
	}), nil
}
