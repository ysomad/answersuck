package round

import (
	"context"

	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) AddTopic(ctx context.Context, r *pb.AddTopicRequest) (*emptypb.Empty, error) {
	return nil, nil
}
