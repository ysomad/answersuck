package v1

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) CreatePlayer(ctx context.Context, p *pb.CreatePlayerRequest) (*emptypb.Empty, error) {
	if err := p.Validate(); err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	if err := h.player.Create(ctx, p.Nickname, p.Email, p.Password); err != nil {
		// TODO:
		// 1. Handle email already exists
		// 2. Handle nickname already exists
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return new(emptypb.Empty), nil
}
