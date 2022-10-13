package email

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
)

func (s *server) VerifyEmail(ctx context.Context, r *pb.VerifyEmailRequest) (*pb.Account, error) {
	return &pb.Account{}, nil
}
