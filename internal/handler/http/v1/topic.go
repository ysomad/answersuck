package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/pkg/logging"
)

type topicService interface {
	Create(ctx context.Context, req dto.TopicCreateRequest) (domain.Topic, error)
	GetAll(ctx context.Context) ([]*domain.Topic, error)
}

type topicHandler struct {
	t       errorTranslator
	log     logging.Logger
	service topicService
}

func newTopicHandler(r *gin.RouterGroup, d *Deps) {
	h := &topicHandler{
		t:       d.ErrorTranslator,
		log:     d.Logger,
		service: d.TopicService,
	}

	topics := r.Group("topics")
	{
		topics.GET("", h.getAll)
	}

	authenticated := topics.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
	{
		authenticated.POST("", h.create)
	}
}

func (h *topicHandler) create(c *gin.Context) {
	var r dto.TopicCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, err := h.service.Create(c.Request.Context(), r)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - topic - create - h.service.Create: %w", err))

		if errors.Is(err, domain.ErrLanguageNotFound) {
			abortWithError(c, http.StatusBadRequest, domain.ErrLanguageNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *topicHandler) getAll(c *gin.Context) {
	t, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - topic - getAll - h.service.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, t)
}
