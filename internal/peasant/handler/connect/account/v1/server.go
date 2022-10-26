package v1

import (
	"context"

	"github.com/ysomad/answersuck/rpc/peasant/v1/v1connect"

	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"

	"github.com/ysomad/answersuck/logger"
)

var _ v1connect.AccountServiceHandler = &server{}

type accountService interface {
	Create(ctx context.Context, args dto.AccountCreateArgs) (*domain.Account, error)
	GetByID(ctx context.Context, accountID string) (*domain.Account, error)
	DeleteByID(ctx context.Context, accountID, password string) error
}

type server struct {
	log            logger.Logger
	accountService accountService
}

func NewServer(l logger.Logger, s accountService) *server {
	return &server{
		log:            l,
		accountService: s,
	}
}
