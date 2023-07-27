package v1

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) CreatePlayer(ctx context.Context, r *pb.CreatePlayerRequest) (*emptypb.Empty, error) {
	if r.Email == "" {
		return nil, twirp.RequiredArgumentError("email")
	}

	if r.Nickname == "" {
		return nil, twirp.RequiredArgumentError("nickname")
	}

	if r.Password == "" {
		return nil, twirp.RequiredArgumentError("password")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	if err := h.player.Create(ctx, r.Nickname, r.Email, r.Password); err != nil {
		if errors.Is(err, apperr.PlayerAlreadyExists) {
			return nil, twirp.AlreadyExists.Error(err.Error())
		}

		return nil, twirp.InternalError(err.Error())
	}

	return new(emptypb.Empty), nil
}
