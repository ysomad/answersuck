package round

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
	_ pb.RoundService  = &Handler{}
)

type UseCase interface {
	Create(ctx context.Context, r entity.Round) (roundID int32, err error)
	Update(ctx context.Context, r entity.Round) error
	GetAll(ctx context.Context, packID int32) ([]entity.Round, error)
}

type sessionManager interface {
	Get(context.Context, string) (*session.Session, error)
}

type Handler struct {
	round   UseCase
	session sessionManager
}

func NewHandler(uc UseCase, sm sessionManager) *Handler {
	return &Handler{
		round:   uc,
		session: sm,
	}
}

func (h *Handler) Handle(m *http.ServeMux) {
	s := pb.NewRoundServiceServer(h,
		twirp.WithServerHooks(hooks.NewLogging()),
		twirp.WithServerHooks(hooks.NewSession(h.session)),
	)
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}
