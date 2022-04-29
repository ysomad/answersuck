package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	repository "github.com/answersuck/vault/internal/repository/psql"

	"github.com/answersuck/vault/pkg/logging"
)

type sessionService interface {
	GetById(ctx context.Context, sid string) (*domain.Session, error)
	GetAll(ctx context.Context, aid string) ([]*domain.Session, error)
	Terminate(ctx context.Context, sid string) error
	TerminateWithExcept(ctx context.Context, aid, sid string) error
}

type sessionHandler struct {
	t       errorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	service sessionService
}

func newSessionHandler(handler *gin.RouterGroup, d *Deps) {
	h := &sessionHandler{
		t:       d.GinTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.SessionService,
	}

	g := handler.Group("sessions")
	{
		authenticated := g.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
		{
			secure := authenticated.Group("", tokenMiddleware(d.Logger, d.AuthService))
			{
				secure.DELETE(":sessionId", h.terminate)
				secure.DELETE("", h.terminateAll)
			}

			authenticated.GET("", h.get)
		}
	}
}

func (h *sessionHandler) get(c *gin.Context) {
	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessions, err := h.service.GetAll(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get - h.service.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *sessionHandler) terminate(c *gin.Context) {
	currSid := getSessionId(c)

	sid := c.Param("sessionId")
	if currSid == sid {
		abortWithError(c, http.StatusBadRequest, domain.ErrSessionCannotBeTerminated, "")
		return
	}

	if err := h.service.Terminate(c.Request.Context(), sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - terminate - h.service.Terminate: %w", err))

		if errors.Is(err, repository.ErrNoAffectedRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - GetAccountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	if err = h.service.TerminateWithExcept(c.Request.Context(), aid, getSessionId(c)); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.service.TerminateWithExcept: %w", err)

		if errors.Is(err, repository.ErrNoAffectedRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
