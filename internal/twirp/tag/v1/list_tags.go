package v1

import (
	"context"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/tag/v1"
)

func (h *Handler) ListTags(ctx context.Context, r *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, err.Error())
	}

	return nil, nil
}
