package round

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
)

func (h *Handler) ListRounds(ctx context.Context, r *pb.ListRoundsRequest) (*pb.ListRoundsResponse, error) {
	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	rounds, err := h.round.GetAll(ctx, r.PackId)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	res := make([]*pb.Round, len(rounds))

	for i, round := range rounds {
		res[i] = &pb.Round{
			Id:       round.ID,
			Name:     round.Name,
			Position: int32(round.Position),
			PackId:   round.PackID,
		}
	}

	return &pb.ListRoundsResponse{Rounds: res}, nil
}
