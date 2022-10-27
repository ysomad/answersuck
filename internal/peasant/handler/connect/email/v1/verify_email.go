package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) VerifyEmail(ctx context.Context, r *connect.Request[pb.VerifyEmailRequest]) (*connect.Response[pb.VerifyEmailResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.emailService.Verify(ctx, r.Msg.GetToken())
	if err != nil {
		s.log.Error(err.Error())

		switch {
		case errors.Is(err, domain.ErrEmailAlreadyVerified):
			return nil, connect.NewError(connect.CodeInvalidArgument, domain.ErrEmailAlreadyVerified)
		case errors.Is(err, domain.ErrEmailVerifTokenExpired):
			return nil, connect.NewError(connect.CodePermissionDenied, domain.ErrEmailVerifTokenExpired)
		}

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
