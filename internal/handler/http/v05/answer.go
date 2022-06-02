package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/answer"
	"github.com/answersuck/vault/internal/domain/media"

	"github.com/answersuck/vault/pkg/logging"
)

type AnswerService interface {
	Create(ctx context.Context, r answer.CreateRequest) (answer.Answer, error)
}

type answerHandler struct {
	t   ErrorTranslator
	cfg *config.Aggregate
	log logging.Logger

	service AnswerService
}

func newAnswerHandler(r *gin.RouterGroup, d *Deps) {
	h := &answerHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.AnswerService,
	}

	answers := r.Group("answers",
		protectionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
	{
		answers.POST("", h.create)
	}
}

func (h *answerHandler) create(c *gin.Context) {
	var r answer.CreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	a, err := h.service.Create(c.Request.Context(), r)
	if err != nil {
		h.log.Error("http - v1 - answer - create - h.service.Create: %w", err)

		switch {
		case errors.Is(err, answer.ErrMimeTypeNotAllowed):
			abortWithError(c, http.StatusBadRequest, answer.ErrMimeTypeNotAllowed, "")
			return
		case errors.Is(err, media.ErrNotFound):
			abortWithError(c, http.StatusBadRequest, answer.ErrMediaNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, a)
}
