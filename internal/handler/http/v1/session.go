package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

type SessionService interface {
	GetByIdWithVerified(ctx context.Context, sessionId string) (*session.SessionWithVerified, error)
	GetById(ctx context.Context, sessionId string) (*session.Session, error)
	GetAll(ctx context.Context, accountId string) ([]*session.Session, error)
	Terminate(ctx context.Context, sessionId string) error
	TerminateWithExcept(ctx context.Context, accountId, sessionId string) error
}

type sessionHandler struct {
	t       ErrorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	service SessionService
}

func newSessionHandler(r *gin.RouterGroup, d *Deps) {
	h := &sessionHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.SessionService,
	}

	sessions := r.Group("sessions")

	authenticated := sessions.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
	{
		authenticated.GET("", h.getAll)
	}

	secure := authenticated.Group("", tokenMiddleware(d.Logger, d.AuthService))
	{
		secure.DELETE(":sessionId", h.terminate)
		secure.DELETE("", h.terminateAll)
	}
}

func (h *sessionHandler) getAll(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - getAccountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessions, err := h.service.GetAll(c.Request.Context(), accountId)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - h.service.GetAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

const (
	paramSessionId = "sessionId"
)

func (h *sessionHandler) terminate(c *gin.Context) {
	currSessionId := getSessionId(c)

	sessionId := c.Param(paramSessionId)
	if currSessionId == sessionId {
		abortWithError(c, http.StatusBadRequest, session.ErrCannotBeTerminated, "")
		return
	}

	if err := h.service.Terminate(c.Request.Context(), sessionId); err != nil {
		h.log.Error("http - v1 - session - terminate - h.service.Terminate: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - getAccountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	if err = h.service.TerminateWithExcept(c.Request.Context(), accountId, getSessionId(c)); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.service.TerminateWithExcept: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
