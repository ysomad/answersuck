package v1

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) VerifyEmail(ctx context.Context, r *connect.Request[pb.VerifyEmailRequest]) (*connect.Response[pb.VerifyEmailResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.emailService.Verify(ctx, r.Msg.GetVerificationCode())
	if err != nil {
		s.log.Error(err.Error())

		// TODO: handle specific errors

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.VerifyEmailResponse{
		Account: &pb.Account{
			Id:            a.ID,
			Email:         a.Email,
			Username:      a.Username,
			EmailVerified: a.EmailVerified,
			Archived:      a.Archived,
			CreationTime:  timestamppb.New(a.CreatedAt),
			UpdateTime:    timestamppb.New(a.UpdatedAt),
		},
	}), nil
}
