package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/user/domain"
	v1 "github.com/ysomad/answersuck/rpc/user/account/v1"
)

func (s *server) CreateAccount(ctx context.Context, r *connect.Request[v1.CreateAccountRequest]) (*connect.Response[v1.CreateAccountResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.service.Create(ctx, r.Msg.Email, r.Msg.Username, r.Msg.Password)
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
		&v1.CreateAccountResponse{
			Account: &v1.Account{
				Id:            a.ID,
				Email:         a.Email,
				Username:      a.Username,
				EmailVerified: a.EmailVerified,
				CreationTime:  t,
				UpdateTime:    t,
			}},
	), nil
}
