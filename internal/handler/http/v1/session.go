package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/service/repository"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
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
	aid, err := GetAccountId(c)
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
	currSid := GetSessionId(c)

	sid := c.Param("sessionId")
	if currSid == sid {
		abortWithError(c, http.StatusBadRequest, domain.ErrSessionCannotBeTerminated, "")
		return
	}

	if err := h.session.Terminate(c.Request.Context(), sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - session - terminate - h.session.Terminate: %w", err))

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
	aid, err := GetAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - accountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return

	}

	currSid := GetSessionId(c)

	if err = h.session.TerminateWithExcept(c.Request.Context(), aid, currSid); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.session.TerminateAll: %w", err)

		if errors.Is(err, repository.ErrNoAffectedRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
