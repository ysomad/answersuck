package v1

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

type SessionService interface {
	GetByIdWithVerified(ctx context.Context, sessionId string) (*session.WithAccountDetails, error)
	GetById(ctx context.Context, sessionId string) (*session.Session, error)
	GetAll(ctx context.Context, accountId string) ([]*session.Session, error)
	Terminate(ctx context.Context, sessionId string) error
	TerminateWithExcept(ctx context.Context, accountId, sessionId string) error
}

type sessionHandler struct {
	cfg     *config.Aggregate
	log     logging.Logger
	v       ValidationModule
	service SessionService
}

func newSessionRouter(d *Deps) *fiber.App {
	h := &sessionHandler{
		cfg:     d.Config,
		log:     d.Logger,
		v:       d.ValidationModule,
		service: d.SessionService,
	}

	r := fiber.New()

	authenticated := r.Group("/", sessionMW(d.Logger, &d.Config.Session, d.SessionService))
	{
		authenticated.Get("/", h.getAll)
	}

	protected := authenticated.Group("/", tokenMW(d.Logger, d.AuthService))
	{
		protected.Delete(":sessionId", h.terminate)
		protected.Delete("/", h.terminateAll)
	}

	return r
}

func (h *sessionHandler) getAll(c *fiber.Ctx) error {

	return c.SendStatus(http.StatusOK)
}

func (h *sessionHandler) terminate(c *fiber.Ctx) error {

	return c.SendStatus(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *fiber.Ctx) error {

	return c.SendStatus(http.StatusNoContent)
}
