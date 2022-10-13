package email

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) SendEmailVerification(ctx context.Context, r *pb.SendEmailVerificationRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
