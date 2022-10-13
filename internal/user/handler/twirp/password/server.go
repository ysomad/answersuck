package password

import pb "github.com/ysomad/answersuck/rpc/user"

var _ pb.PasswordService = &server{}

type server struct{}

func NewServer() *server {
	return &server{}
}
