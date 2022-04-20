package v1

import (
	"fmt"
	"github.com/answersuck/vault/internal/service"
	"github.com/answersuck/vault/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type languageHandler struct {
	log      logging.Logger
	language service.Language
}

func newLanguageHandler(handler *gin.RouterGroup, d *Deps) {
	h := &languageHandler{
		log:      d.Logger,
		language: d.LanguageService,
	}

	g := handler.Group("languages")
	{
		g.GET("", h.getAll)
	}
}

func (h *languageHandler) getAll(c *gin.Context) {
	l, err := h.language.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - language - getAll - h.language.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, l)
}
