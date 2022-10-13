package email

import pb "github.com/ysomad/answersuck/rpc/user"

var _ pb.EmailService = &server{}

type server struct{}

func NewServer() *server {
	return &server{}
}
