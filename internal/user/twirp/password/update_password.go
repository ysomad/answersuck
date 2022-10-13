package password

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) UpdatePassword(ctx context.Context, r *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
