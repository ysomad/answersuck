package app

import (
	"github.com/ysomad/answersuck/internal/user/twirp/account"
	"github.com/ysomad/answersuck/internal/user/twirp/email"
	"github.com/ysomad/answersuck/internal/user/twirp/password"
	pb "github.com/ysomad/answersuck/rpc/user"
)

func newAccountServiceServer() pb.TwirpServer {
	s := account.NewServer()
	return pb.NewAccountServiceServer(s)
}

func newPasswordServiceServer() pb.TwirpServer {
	s := password.NewServer()
	return pb.NewPasswordServiceServer(s)
}

func newEmailServiceServer() pb.TwirpServer {
	s := email.NewServer()
	return pb.NewEmailServiceServer(s)
}
