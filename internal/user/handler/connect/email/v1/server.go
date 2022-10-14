package v1

import (
	"github.com/ysomad/answersuck/rpc/user/email/v1/emailv1connect"

	"github.com/ysomad/answersuck/logger"
)

var _ emailv1connect.EmailServiceHandler = &server{}

type server struct {
	log logger.Logger
}

func NewServer(l logger.Logger) *server {
	return &server{
		log: l,
	}
}
