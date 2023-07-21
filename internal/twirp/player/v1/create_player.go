package v1

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) CreatePlayer(ctx context.Context, p *pb.CreatePlayerRequest) (*emptypb.Empty, error) {
	if err := p.Validate(); err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	if err := h.player.Create(ctx, p.Nickname, p.Email, p.Password); err != nil {
		if errors.Is(err, apperr.ErrPlayerAlreadyExist) {
			return nil, twirp.AlreadyExists.Error(err.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return new(emptypb.Empty), nil
}
