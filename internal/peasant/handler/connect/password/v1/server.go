package v1

import (
	"github.com/ysomad/answersuck/rpc/peasant/v1/v1connect"

	"github.com/ysomad/answersuck/logger"
)

var _ v1connect.PasswordServiceHandler = &server{}

type server struct {
	log logger.Logger
}

func NewServer(l logger.Logger) *server {
	return &server{
		log: l,
	}
}
