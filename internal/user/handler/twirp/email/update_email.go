package email

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
)

func (s *server) UpdateEmail(ctx context.Context, r *pb.UpdateEmailRequest) (*pb.Account, error) {
	return &pb.Account{}, nil
}
