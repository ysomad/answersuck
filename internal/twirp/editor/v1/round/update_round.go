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

func (h *Handler) UpdateRound(ctx context.Context, r *pb.UpdateRoundRequest) (*pb.UpdateRoundResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if r.RoundId == 0 {
		return nil, twirp.RequiredArgumentError("round_id")
	}

	if r.RoundName == "" {
		return nil, twirp.RequiredArgumentError("round_name")
	}

	if r.RoundPosition == 0 {
		return nil, twirp.RequiredArgumentError("round_position")
	}

	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	round := entity.Round{
		ID:       r.RoundId,
		Name:     r.RoundName,
		PackID:   r.PackId,
		Position: int16(r.RoundPosition),
	}

	if err := h.round.Update(ctx, session.User.ID, round); err != nil {
		switch {
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		case errors.Is(err, apperr.RoundNotFound):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundNotFound)
		case errors.Is(err, apperr.PackNotFound):
			return nil, twirp.InvalidArgument.Error(apperr.MsgPackNotFound)
		}
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.UpdateRoundResponse{
		Round: &pb.Round{
			Id:       round.ID,
			Name:     round.Name,
			Position: int32(round.Position),
			PackId:   round.PackID,
		},
	}, nil
}
