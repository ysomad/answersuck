package v1

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetPlayer(ctx context.Context, p *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	player, err := h.player.GetOne(ctx, p.Nickname)
	if err != nil {
		// TODO:
		// 1. Handle player not found
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}

	return &pb.GetPlayerResponse{
		Player: &pb.Player{
			Nickname:      player.Nickname,
			Email:         player.Email,
			DisplayName:   player.DisplayName,
			EmailVerified: player.EmailVerified,
			CreateTime:    timestamppb.New(player.CreatedAt),
		},
	}, nil
}
