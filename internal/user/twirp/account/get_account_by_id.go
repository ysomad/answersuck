package account

import (
	"context"
	pb "github.com/ysomad/answersuck/rpc/user"
)

func (s *server) GetAccountById(ctx context.Context, r *pb.GetAccountByIdRequest) (*pb.Account, error) {
	return &pb.Account{}, nil
}
