package account

import (
	"context"

	"github.com/ysomad/answersuck/user/internal/gen/proto/account/accountconnect"

	"github.com/ysomad/answersuck/user/internal/entity"
	"github.com/ysomad/answersuck/user/internal/service/dto"

	"github.com/ysomad/answersuck/pkg/logger"
)

var _ accountconnect.AccountServiceHandler = &server{}

type accountService interface {
	Create(ctx context.Context, p dto.AccountCreateParams) (*entity.Account, error)
	GetByID(ctx context.Context, accountID string) (*entity.Account, error)
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
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
