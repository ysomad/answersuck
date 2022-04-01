package v1

import (
	"fmt"
	"github.com/answersuck/vault/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/service"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/validation"
)

type sessionHandler struct {
	t       validation.ErrorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	session service.Session
}

func newSessionHandler(handler *gin.RouterGroup, d *Deps) {
	h := &sessionHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		session: d.SessionService,
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
	aid, err := accountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get - accountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessions, err := h.session.GetAll(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get - h.session.GetAll: %w", err))

		// TODO: handle specific errors

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *sessionHandler) terminate(c *gin.Context) {
	if err := h.session.Terminate(c.Request.Context(), c.Param("sessionID")); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - sessionService - terminate - h.sessionService.Terminate: %w", err))

		// TODO: handle specific errors

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	sid, err := sessionId(c)
	if err != nil {
		h.log.Error("http - v1 - sessionService - terminateAll - sessionId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	aid, err := accountId(c)
	if err != nil {
		h.log.Error("http - v1 - sessionService - terminateAll - accountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	if err = h.session.TerminateAll(c.Request.Context(), aid, sid); err != nil {
		h.log.Error("http - v1 - sessionService - terminateAll - h.session.TerminateAll: %w", err)

		// TODO: handle specific errors

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
