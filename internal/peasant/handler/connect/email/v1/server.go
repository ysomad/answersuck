package v1

import (
	"context"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	"github.com/ysomad/answersuck/rpc/peasant/v1/v1connect"

	"github.com/ysomad/answersuck/logger"
)

var _ v1connect.EmailServiceHandler = &server{}

type emailService interface {
	Update(ctx context.Context, args dto.UpdateEmailArgs) (*domain.Account, error)
	Verify(ctx context.Context, code string) (*domain.Account, error)
	SendVerification(ctx context.Context, accountID string) error
}

type server struct {
	log          logger.Logger
	emailService emailService
}

func NewServer(l logger.Logger, s emailService) *server {
	return &server{
		log:          l,
		emailService: s,
	}
}
