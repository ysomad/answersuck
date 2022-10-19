package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) SetPassword(ctx context.Context, r *connect.Request[pb.SetPasswordRequest]) (*connect.Response[pb.SetPasswordResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.passwordService.Set(ctx, r.Msg.GetToken(), r.Msg.GetNewPassword())
	if err != nil {
		s.log.Error(err.Error())

		if errors.Is(err, domain.ErrPasswordTokenExpired) {
			return nil, connect.NewError(connect.CodePermissionDenied, domain.ErrPasswordTokenExpired)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.SetPasswordResponse{
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
