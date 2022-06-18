package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain/tag"

	"github.com/answersuck/vault/pkg/logging"
)

type TagService interface {
	CreateMultiple(ctx context.Context, r []tag.CreateReq) ([]*tag.Tag, error)
	GetAll(ctx context.Context) ([]*tag.Tag, error)
}

type tagHandler struct {
	t       ErrorTranslator
	log     logging.Logger
	service TagService
}

func newTagHandler(r *gin.RouterGroup, d *Deps) {
	h := &tagHandler{
		t:       d.ErrorTranslator,
		log:     d.Logger,
		service: d.TagService,
	}

	tags := r.Group("tags")
	{
		tags.GET("", h.getAll)
		tags.POST("",
			protectionMiddleware(d.Logger, &d.Config.Session, d.SessionService),
			h.createMultiple)
	}
}

func (h *tagHandler) createMultiple(c *gin.Context) {
	var r tag.CreateMultipleReq

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, err := h.service.CreateMultiple(c.Request.Context(), r.Tags)
	if err != nil {
		h.log.Error("http - v1 - tag - createMultiple - h.service.CreateMultiple: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, listResponse{Result: t})
}

func (h *tagHandler) getAll(c *gin.Context) {
	t, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error("http - v1 - tag - getAll - h.service.GetAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, listResponse{Result: t})
}
