package v1

import (
	"context"

	"github.com/ysomad/answersuck/rpc/user/account/v1/accountv1connect"

	"github.com/ysomad/answersuck/internal/user/domain"

	"github.com/ysomad/answersuck/logger"
)

var _ accountv1connect.AccountServiceHandler = &server{}

type accountService interface {
	Create(ctx context.Context, email, username, plainPassword string) (*domain.Account, error)
	GetByID(ctx context.Context, accountID string) (*domain.Account, error)
	GetByEmail(ctx context.Context, email string) (*domain.Account, error)
	DeleteByID(ctx context.Context, accountID string) error
}

type server struct {
	log     logger.Logger
	service accountService
}

func NewServer(l logger.Logger, s accountService) *server {
	return &server{
		log:     l,
		service: s,
	}
}
