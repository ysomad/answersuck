package account

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/user/internal/gen/proto/account"
	"github.com/ysomad/answersuck/user/internal/service/dto"
)

func (s *server) CreateAccount(
	ctx context.Context,
	r *connect.Request[account.CreateAccountRequest],
) (*connect.Response[account.CreateAccountResponse], error) {
	if err := r.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	a, err := s.service.Create(ctx, dto.AccountCreateParams{
		Email:         r.Msg.GetEmail(),
		Username:      r.Msg.GetUsername(),
		PlainPassword: r.Msg.GetPassword(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &account.CreateAccountResponse{Account: &account.Account{
		Id:          a.ID,
		Email:       a.Email,
		Username:    a.Username,
		CreatedTime: timestamppb.New(a.CreatedAt),
		UpdatedTime: timestamppb.New(a.UpdatedAt),
	}}

	return connect.NewResponse[account.CreateAccountResponse](resp), nil
}
