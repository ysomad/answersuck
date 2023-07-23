package v1

import (
	"context"
	"net/http"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/player/v1"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
)

var (
	_ apptwirp.Handler = &Handler{}
	_ pb.PlayerService = &Handler{}
)

type UseCase interface {
	Create(ctx context.Context, nickname, email, password string) error
	Get(ctx context.Context, login string) (*entity.Player, error)
}

type Handler struct {
	player UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{
		player: uc,
	}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewPlayerServiceServer(h, hooks.NewLogging())
	m.Handle(s.PathPrefix(), s)
}
