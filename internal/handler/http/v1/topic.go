package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
)

type topicService interface {
	GetAll(ctx context.Context) ([]*domain.Topic, error)
}

type topicHandler struct {
	log     logging.Logger
	service topicService
}

func newTopicHandler(handler *gin.RouterGroup, d *Deps) {
	h := &topicHandler{
		log:     d.Logger,
		service: d.TopicService,
	}

	g := handler.Group("topics")
	{
		g.GET("", h.getAll)
	}
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
