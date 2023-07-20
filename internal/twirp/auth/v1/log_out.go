package v1

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) LogOut(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
