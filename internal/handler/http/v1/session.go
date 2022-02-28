package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/internal/service"

	"github.com/quizlyfun/quizly-backend/pkg/logging"
	"github.com/quizlyfun/quizly-backend/pkg/validation"
)

type sessionHandler struct {
	validation.ErrorTranslator
	log     logging.Logger
	session service.Session
}

func newSessionHandler(handler *gin.RouterGroup, d *Deps) {
	h := &sessionHandler{
		d.ErrorTranslator,
		d.Logger,
		d.SessionService,
	}

	g := handler.Group("sessions")
	{
		authenticated := g.Group("", sessionMiddleware(d.Logger, d.Config, d.SessionService))
		{
			secure := authenticated.Group("", tokenMiddleware(d.Logger, d.AuthService))
			{
				secure.DELETE(":sessionID", h.terminate)
				secure.DELETE("", h.terminateAll)
			}

			authenticated.GET("", h.get)
		}
	}
}

func (h *sessionHandler) get(c *gin.Context) {
	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get - accountID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessions, err := h.session.GetAll(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - get: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *sessionHandler) terminate(c *gin.Context) {
	currSid, err := sessionID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminate - sessionID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := h.session.Terminate(c.Request.Context(), c.Param("sessionID"), currSid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - terminate - h.session.Terminate: %w", err))

		if errors.Is(err, domain.ErrSessionNotTerminated) {
			abortWithError(c, http.StatusBadRequest, domain.ErrSessionNotTerminated, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(c *gin.Context) {
	currSid, err := sessionID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - sessionID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	aid, err := accountID(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - accountID: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	if err := h.session.TerminateAll(c.Request.Context(), aid, currSid); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.session.TerminateAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
