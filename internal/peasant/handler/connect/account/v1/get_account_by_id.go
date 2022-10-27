package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/peasant/domain"

	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) GetAccountById(ctx context.Context, r *connect.Request[pb.GetAccountByIdRequest]) (*connect.Response[pb.GetAccountByIdResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.accountService.GetByID(ctx, r.Msg.GetAccountId())
	if err != nil {
		s.log.Error(err.Error())

		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, domain.ErrAccountNotFound)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(
		&pb.GetAccountByIdResponse{
			Account: &pb.Account{
				Id:            a.ID,
				Email:         a.Email,
				Username:      a.Username,
				EmailVerified: a.EmailVerified,
				Archived:      a.Archived,
				CreationTime:  timestamppb.New(a.CreatedAt),
				UpdateTime:    timestamppb.New(a.UpdatedAt),
			},
		},
	), nil
}
