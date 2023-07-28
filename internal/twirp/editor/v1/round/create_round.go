package round

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/twirp/common"
)

func (h *Handler) CreateRound(ctx context.Context, r *pb.CreateRoundRequest) (*pb.CreateRoundResponse, error) {
	var err error

	if _, err = common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	if r.RoundName == "" {
		return nil, twirp.RequiredArgumentError("round_name")
	}

	if r.RoundPosition == 0 {
		return nil, twirp.RequiredArgumentError("round_position")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	round := entity.Round{
		Name:     r.RoundName,
		PackID:   r.PackId,
		Position: int16(r.RoundPosition),
	}

	round.ID, err = h.round.Create(ctx, round)
	if err != nil {
		switch {
		case errors.Is(err, apperr.PackNotFound):
			return nil, twirp.InvalidArgumentError("pack_id", apperr.MsgPackNotFound)
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateRoundResponse{
		Round: &pb.Round{
			Id:       round.ID,
			Name:     round.Name,
			Position: int32(round.Position),
			PackId:   round.PackID,
		},
	}, nil
}
