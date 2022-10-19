package v1

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) SendVerification(ctx context.Context, r *connect.Request[pb.SendVerificationRequest]) (*connect.Response[pb.SendVerificationResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	v, err := s.emailService.CreateVerification(ctx, r.Msg.GetAccountId())
	if err != nil {
		s.log.Error(err.Error())

		// TODO: connect SendVerification handle specific errors

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// TODO: send email verification

	return connect.NewResponse(&pb.SendVerificationResponse{
		EmailVerification: &pb.EmailVerification{
			AccountId:        v.AccountID,
			VerificationCode: v.Code,
			ExpirationTime:   timestamppb.New(v.ExpiresAt),
		}}), nil
}
