package v1

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	"github.com/ysomad/answersuck/rpc/peasant/v1/v1connect"

	"github.com/ysomad/answersuck/logger"
)

var _ v1connect.EmailServiceHandler = &server{}

type accountEmailService interface {
	Update(ctx context.Context, args dto.UpdateEmailArgs) (*domain.Account, error)
	SendVerification(ctx context.Context, accountID string) error
	Verify(ctx context.Context, code string) (*domain.Account, error)
}

type server struct {
	log          logger.Logger
	accountEmail accountEmailService
}

func NewServer(l logger.Logger, s accountEmailService) *server {
	return &server{
		log:          l,
		accountEmail: s,
	}
}
