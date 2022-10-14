package v1

import (
	"github.com/ysomad/answersuck/rpc/user/password/v1/passwordv1connect"

	"github.com/ysomad/answersuck/logger"
)

var _ passwordv1connect.PasswordServiceHandler = &server{}

type server struct {
	log logger.Logger

	passwordv1connect.UnimplementedPasswordServiceHandler
}

func NewServer(l logger.Logger) *server {
	return &server{
		log: l,
	}
}
