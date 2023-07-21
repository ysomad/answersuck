package v1

import (
	"context"
	"net/http"

	pb "github.com/ysomad/answersuck/internal/gen/api/auth/v1"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/session"
	"github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ twirp.Handler  = &Handler{}
	_ pb.AuthService = &Handler{}
)

type UseCase interface {
	LogIn(ctx context.Context, login, password string, fp appctx.FootPrint) (*session.Session, error)
}

type Handler struct {
	auth UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{auth: uc}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewAuthServiceServer(h)
	m.Handle(s.PathPrefix(), middleware.WithFootPrint(middleware.WithSessionID(s)))
}
