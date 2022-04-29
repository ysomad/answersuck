package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
)

type languageService interface {
	GetAll(ctx context.Context) ([]*domain.Language, error)
}

type languageHandler struct {
	log     logging.Logger
	service languageService
}

func newLanguageHandler(handler *gin.RouterGroup, d *Deps) {
	h := &languageHandler{
		log:     d.Logger,
		service: d.LanguageService,
	}

	g := handler.Group("languages")
	{
		g.GET("", h.getAll)
	}
}

func (h *languageHandler) getAll(c *gin.Context) {
	l, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - language - getAll - h.service.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, l)
}
