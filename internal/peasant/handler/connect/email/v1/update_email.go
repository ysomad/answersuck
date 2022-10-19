package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) UpdateEmail(ctx context.Context, r *connect.Request[pb.UpdateEmailRequest]) (*connect.Response[pb.UpdateEmailResponse], error) {
	if err := r.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.emailService.Update(ctx, dto.UpdateEmailArgs{
		AccountID:     r.Msg.GetAccountId(),
		NewEmail:      r.Msg.GetNewEmail(),
		PlainPassword: r.Msg.GetPassword(),
	})
	if err != nil {
		s.log.Error(err.Error())

		switch {
		case errors.Is(err, domain.ErrIncorrectPassword):
			return nil, connect.NewError(connect.CodePermissionDenied, domain.ErrIncorrectPassword)
		case errors.Is(err, domain.ErrEmailTaken):
			return nil, connect.NewError(connect.CodeAlreadyExists, domain.ErrEmailTaken)
		case errors.Is(err, domain.ErrAccountNotFound):
			return nil, connect.NewError(connect.CodeNotFound, domain.ErrAccountNotFound)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateEmailResponse{
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
