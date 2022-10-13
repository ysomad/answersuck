package account

import (
	pb "github.com/ysomad/answersuck/rpc/user"
)

var _ pb.AccountService = &server{}

type server struct{}

func NewServer() *server {
	return &server{}
}
