package pack

import (
	"context"
	"net/http"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ apptwirp.Handler = &Handler{}
	_ pb.PackService   = &Handler{}
)

type UseCase interface {
	Save(ctx context.Context, p *entity.Pack, tags []string) (packID int32, err error)
	GetOne(ctx context.Context, packID int32) (*entity.Pack, error)
	GetTags(ctx context.Context, packID int32) ([]string, error)
}

type sessionManager interface {
	Get(context.Context, string) (*session.Session, error)
}

type Handler struct {
	pack    UseCase
	session sessionManager
}

func NewHandler(uc UseCase, sm sessionManager) *Handler {
	return &Handler{
		pack:    uc,
		session: sm,
	}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewPackServiceServer(h,
		twirp.WithServerHooks(hooks.NewLogging()),
		twirp.WithServerHooks(hooks.NewSession(h.session)),
	)
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}