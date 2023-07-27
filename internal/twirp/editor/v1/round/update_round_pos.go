package round

import (
	"context"

	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) UpdateRoundPos(ctx context.Context, r *pb.UpdateRoundPosRequest) (*emptypb.Empty, error) {
	return nil, nil
}
