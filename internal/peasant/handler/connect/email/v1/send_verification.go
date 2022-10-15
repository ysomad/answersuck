package v1

import (
	"context"

	"github.com/bufbuild/connect-go"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) SendVerification(ctx context.Context, r *connect.Request[pb.SendVerificationRequest]) (*connect.Response[pb.SendVerificationResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := s.emailService.SendVerification(ctx, r.Msg.GetAccountId()); err != nil {
		s.log.Error(err.Error())

		// TODO: handle specific errors

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.SendVerificationResponse{}), nil
}
