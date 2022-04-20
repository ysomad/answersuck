package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/service"
	"github.com/answersuck/vault/pkg/logging"
)

type topicHandler struct {
	log   logging.Logger
	topic service.Topic
}

func newTopicHandler(handler *gin.RouterGroup, d *Deps) {
	h := &topicHandler{
		log:   d.Logger,
		topic: d.TopicService,
	}

	g := handler.Group("topics")
	{
		g.GET("", h.getAll)
	}
}

func (h *topicHandler) getAll(c *gin.Context) {
	t, err := h.topic.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - topic - getAll - h.topic.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, t)
}
