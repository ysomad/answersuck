package v1

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) UpdatePassword(ctx context.Context, r *connect.Request[pb.UpdatePasswordRequest]) (*connect.Response[pb.UpdatePasswordResponse], error) {
	if err := r.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.passwordService.Update(ctx, dto.UpdatePasswordArgs{
		AccountID:   r.Msg.GetAccountId(),
		OldPassword: r.Msg.GetOldPassword(),
		NewPassword: r.Msg.GetNewPassword(),
	})
	if err != nil {
		s.log.Error(err.Error())

		// TODO: UpdatePassword handle specific errors

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdatePasswordResponse{
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
