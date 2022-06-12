package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

type sessionHandler struct {
	cfg     *config.Aggregate
	log     logging.Logger
	v       ValidationModule
	service SessionService
}

func newSessionHandler(d *Deps) *sessionHandler {
	return &sessionHandler{
		cfg:     d.Config,
		log:     d.Logger,
		v:       d.ValidationModule,
		service: d.SessionService,
	}
}

func newSessionRouter(d *Deps) *fiber.App {
	h := newSessionHandler(d)

	r := fiber.New()

	r.Get("/", sessionMW(d.Logger, &d.Config.Session, d.SessionService), h.getAll)

	requireToken := r.Group("/",
		sessionMW(d.Logger, &d.Config.Session, d.SessionService),
		tokenMW(d.Logger, d.TokenService))
	requireToken.Delete(":sessionId", h.terminate)
	requireToken.Delete("/", h.terminateAll)

	return r
}

func (h *sessionHandler) getAll(c *fiber.Ctx) error {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - getAccountId: %w", err)
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	s, err := h.service.GetAll(c.Context(), accountId)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - h.service.GetAll: %w", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(s)
}

func (h *sessionHandler) terminate(c *fiber.Ctx) error {
	currSessionId, err := getSessionId(c)
	if err != nil {
		h.log.Info("http - v1 - session - terminate - getSessionId: %w", err)
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	sessionId := c.Params("sessionId")
	if currSessionId == sessionId {
		return errorResp(c, fiber.StatusBadRequest, session.ErrCannotBeTerminated, "")
	}

	if err := h.service.Terminate(c.Context(), sessionId); err != nil {
		h.log.Error("http - v1 - session - terminate - h.service.Terminate: %w", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *sessionHandler) terminateAll(c *fiber.Ctx) error {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - getAccountId: %w", err)
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	sessionId, err := getSessionId(c)
	if err != nil {
		h.log.Info("http - v1 - session - terminateAll - getSessionId: %w", err)
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	if err = h.service.TerminateWithExcept(c.Context(), accountId, sessionId); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.service.TerminateWithExcept: %w", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
