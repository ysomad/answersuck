package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain/language"

	"github.com/answersuck/vault/pkg/logging"
)

type LanguageService interface {
	GetAll(ctx context.Context) ([]*language.Language, error)
}

type languageHandler struct {
	log     logging.Logger
	service LanguageService
}

func newLanguageHandler(r *gin.RouterGroup, d *Deps) {
	h := &languageHandler{
		log:     d.Logger,
		service: d.LanguageService,
	}

	languages := r.Group("languages")
	{
		languages.GET("", h.getAll)
	}
}

func (h *languageHandler) getAll(c *gin.Context) {
	l, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error("http - v1 - language - getAll - h.service.GetAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, listResponse{Result: l})
}
