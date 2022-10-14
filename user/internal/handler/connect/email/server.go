package email

import (
	"github.com/ysomad/answersuck/user/internal/gen/proto/email/emailconnect"

	"github.com/ysomad/answersuck/pkg/logger"
)

var _ emailconnect.EmailServiceHandler = &server{}

type server struct {
	log logger.Logger
}

func NewServer(l logger.Logger) *server {
	return &server{
		log: l,
	}
}
