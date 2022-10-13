package password

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) ResetPassword(ctx context.Context, r *pb.ResetPasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
