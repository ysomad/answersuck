package password

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SetPassword(ctx context.Context, r *pb.SetPasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
