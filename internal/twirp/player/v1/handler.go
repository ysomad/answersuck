package v1

import (
	"context"
	"net/http"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
)

var (
	_ apptwirp.Handler = &Handler{}
	_ pb.PlayerService = &Handler{}
)

type UseCase interface {
	Create(ctx context.Context, nickname, email, password string) error
	Get(ctx context.Context, login string) (*entity.Player, error)
}

type sessionManager interface {
	Get(context.Context, string) (*session.Session, error)
}

type Handler struct {
	session sessionManager
	player  UseCase
}

func NewHandler(uc UseCase, s sessionManager) *Handler {
	return &Handler{
		session: s,
		player:  uc,
	}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewPlayerServiceServer(h)
	m.Handle(s.PathPrefix(), s)
}
