package pack

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
)

func (h *Handler) ListPackTags(ctx context.Context, r *pb.ListPackTagsRequest) (*pb.ListPackTagsResponse, error) {
	if r.PackId == 0 {
		return nil, twirp.RequiredArgumentError("pack_id")
	}

	tags, err := h.pack.GetTags(ctx, r.PackId)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.ListPackTagsResponse{Tags: tags}, nil
}
