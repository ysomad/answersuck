package password

import (
	"github.com/ysomad/answersuck/user/internal/gen/proto/password/passwordconnect"

	"github.com/ysomad/answersuck/pkg/logger"
)

var _ passwordconnect.PasswordServiceHandler = &server{}

type server struct {
	log logger.Logger
}

func NewServer(l logger.Logger) *server {
	return &server{
		log: l,
	}
}
