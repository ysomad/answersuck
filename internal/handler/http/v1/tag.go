package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain/tag"

	"github.com/answersuck/vault/pkg/logging"
)

type TagService interface {
	GetAll(ctx context.Context) ([]*tag.Tag, error)
}

type tagHandler struct {
	log     logging.Logger
	service TagService
}

func newTagHandler(r *gin.RouterGroup, d *Deps) {
	h := &tagHandler{
		log:     d.Logger,
		service: d.TagService,
	}

	tags := r.Group("tags")
	{
		tags.GET("", h.getAll)
	}
}

func (h *tagHandler) getAll(c *gin.Context) {
	t, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error("http - v1 - tag - getAll - h.service.GetAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, t)
}
