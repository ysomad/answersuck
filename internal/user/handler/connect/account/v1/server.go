package v1

import (
	"context"

	"github.com/ysomad/answersuck/rpc/user/account/v1/accountv1connect"

	"github.com/ysomad/answersuck/internal/user/entity"
	"github.com/ysomad/answersuck/internal/user/service/dto"

	"github.com/ysomad/answersuck/logger"
)

var _ accountv1connect.AccountServiceHandler = &server{}

type accountService interface {
	Create(ctx context.Context, p dto.AccountCreateParams) (*entity.Account, error)
	GetByID(ctx context.Context, accountID string) (*entity.Account, error)
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
	DeleteByID(ctx context.Context, accountID string) error
}

type server struct {
	log     logger.Logger
	service accountService

	accountv1connect.UnimplementedAccountServiceHandler
}

func NewServer(l logger.Logger, s accountService) *server {
	return &server{
		log:     l,
		service: s,
	}
}
