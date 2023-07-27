package round

import (
	"context"

	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) RemoveTopic(ctx context.Context, r *pb.RemoveTopicRequest) (*emptypb.Empty, error) {
	return nil, nil
}
