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

func (s *server) CreateAccount(ctx context.Context, r *connect.Request[pb.CreateAccountRequest]) (*connect.Response[pb.CreateAccountResponse], error) {
	if err := r.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.accountService.Create(ctx, dto.AccountCreateArgs{
		Email:    r.Msg.GetEmail(),
		Username: r.Msg.GetUsername(),
		Password: r.Msg.GetPassword(),
	})
	if err != nil {
		s.log.Error(err.Error())

		switch {
		case errors.Is(err, domain.ErrEmailTaken):
			return nil, connect.NewError(connect.CodeAlreadyExists, domain.ErrEmailTaken)
		case errors.Is(err, domain.ErrUsernameTaken):
			return nil, connect.NewError(connect.CodeAlreadyExists, domain.ErrUsernameTaken)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// using t as response for creation time and update time
	// because on account create its the same
	t := timestamppb.New(a.CreatedAt)

	return connect.NewResponse(
		&pb.CreateAccountResponse{
			Account: &pb.Account{
				Id:            a.ID,
				Email:         a.Email,
				Username:      a.Username,
				EmailVerified: a.EmailVerified,
				Archived:      a.Archived,
				CreationTime:  t,
				UpdateTime:    t,
			},
		},
	), nil
}
