package account

import (
	"context"

	pb "github.com/ysomad/answersuck/rpc/user"

	"github.com/ysomad/answersuck/internal/user/entity"
	"github.com/ysomad/answersuck/internal/user/service/dto"
	"github.com/ysomad/answersuck/pkg/logger"
)

var _ pb.AccountService = &server{}

type accountService interface {
	Create(ctx context.Context, p dto.AccountCreateParams) (*entity.Account, error)
	GetByID(ctx context.Context, accountID string) (*entity.Account, error)
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
	DeleteByID(ctx context.Context, accountID string) error
}

type server struct {
	log     logger.Logger
	account accountService
}

func NewServer(l logger.Logger, s accountService) *server {
	return &server{
		log:     l,
		account: s,
	}
}
