package auth

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/session"
)

type sessionService interface {
	Create(context.Context, session.User) (*session.Session, error)
	Delete(context.Context, string) error
}

type playerService interface {
	Get(context.Context, string) (*entity.Player, error)
}

type Service struct {
	session sessionService
	player  playerService
}

func NewService(s sessionService, p playerService) *Service {
	return &Service{
		session: s,
		player:  p,
	}
}
