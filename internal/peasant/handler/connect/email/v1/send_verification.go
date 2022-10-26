package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) SendVerification(ctx context.Context, r *connect.Request[pb.SendVerificationRequest]) (*connect.Response[pb.SendVerificationResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	t, err := s.emailService.NotifyWithToken(ctx, r.Msg.GetAccountId())
	if err != nil {
		s.log.Error(err.Error())

		// TODO: connect SendVerification handle specific errors
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, domain.ErrAccountIDNotFound)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.SendVerificationResponse{Token: t.String()}), nil
}
