package v1

import (
	"fmt"
	"github.com/answersuck/vault/internal/service"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type tagHandler struct {
	log logging.Logger
	tag service.Tag
}

func newTagHandler(handler *gin.RouterGroup, d *Deps) {
	h := &tagHandler{
		log: d.Logger,
		tag: d.TagService,
	}

	g := handler.Group("tags")
	{
		g.GET("", h.getAll)
	}
}

func (h *tagHandler) getAll(c *gin.Context) {
	tags, err := h.tag.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - tag - getAll - h.tag.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, tags)
}
