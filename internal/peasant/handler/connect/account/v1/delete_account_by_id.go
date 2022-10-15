package v1

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	pb "github.com/ysomad/answersuck/rpc/peasant/v1"
)

func (s *server) DeleteAccountById(ctx context.Context, r *connect.Request[pb.DeleteAccountByIdRequest]) (*connect.Response[pb.DeleteAccountByIdResponse], error) {
	if err := r.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := s.accountService.DeleteByID(ctx, r.Msg.GetAccountId()); err != nil {
		s.log.Error(err.Error())

		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, domain.ErrAccountNotFound)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteAccountByIdResponse{}), nil
}
