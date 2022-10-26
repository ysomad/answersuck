package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) DeleteAccountById(ctx context.Context, r *connect.Request[pb.DeleteAccountByIdRequest]) (*connect.Response[pb.DeleteAccountByIdResponse], error) {
	if err := r.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := s.accountService.DeleteByID(ctx, r.Msg.GetAccountId(), r.Msg.GetPassword()); err != nil {
		s.log.Error(err.Error())

		switch {
		case errors.Is(err, domain.ErrAccountNotFound):
			return nil, connect.NewError(connect.CodeNotFound, domain.ErrAccountNotFound)
		case errors.Is(err, domain.ErrIncorrectPassword):
			return nil, connect.NewError(connect.CodePermissionDenied, domain.ErrIncorrectPassword)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteAccountByIdResponse{}), nil
}
