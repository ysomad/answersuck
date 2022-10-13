package account

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) DeleteAccount(ctx context.Context, r *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
