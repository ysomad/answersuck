package account

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
)

func (s *server) GetAccountByEmail(ctx context.Context, r *pb.GetAccountByEmailRequest) (*pb.Account, error) {
	return &pb.Account{}, nil
}
