package v1

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	"github.com/ysomad/answersuck/rpc/peasant/v1/v1connect"

	"github.com/ysomad/answersuck/logger"
)

var _ v1connect.PasswordServiceHandler = &server{}

type passwordService interface {
	Update(ctx context.Context, args dto.UpdatePasswordArgs) (*domain.Account, error)
}

type server struct {
	log             logger.Logger
	passwordService passwordService
}

func NewServer(l logger.Logger, s passwordService) *server {
	return &server{
		log:             l,
		passwordService: s,
	}
}
