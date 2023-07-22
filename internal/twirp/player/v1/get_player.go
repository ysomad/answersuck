package v1

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetPlayer(ctx context.Context, p *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	if p.Nickname == "" {
		return nil, twirp.RequiredArgumentError("nickname")
	}

	player, err := h.player.Get(ctx, p.Nickname)
	if err != nil {
		if errors.Is(err, apperr.ErrPlayerNotFound) {
			return nil, twirp.NotFoundError(apperr.ErrPlayerNotFound.Error())
		}

		return nil, twirp.InternalError(err.Error())
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
