package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain/topic"
	"github.com/answersuck/vault/pkg/logging"
)

type TopicService interface {
	Create(ctx context.Context, req topic.CreateRequest) (topic.Topic, error)
	GetAll(ctx context.Context) ([]*topic.Topic, error)
}

type topicHandler struct {
	t       ErrorTranslator
	log     logging.Logger
	service TopicService
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

	protected := topics.Group(
		"",
		sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService),
		protectionMiddleware(d.Logger),
	)
	{
		protected.POST("", h.create)
	}
}

func (h *topicHandler) create(c *gin.Context) {
	var r topic.CreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, err := h.service.Create(c.Request.Context(), r)
	if err != nil {
		h.log.Error("http - v1 - topic - create - h.service.Create: %w", err)

		if errors.Is(err, topic.ErrLanguageNotFound) {
			abortWithError(c, http.StatusBadRequest, topic.ErrLanguageNotFound, "")
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
		h.log.Error("http - v1 - topic - getAll - h.service.GetAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, t)
}
