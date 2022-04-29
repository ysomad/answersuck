package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/pkg/logging"
)

type questionService interface {
	GetAll(ctx context.Context) ([]*domain.Question, error)
}

type questionHandler struct {
	log     logging.Logger
	service questionService
}

func newQuestionHandler(handler *gin.RouterGroup, d *Deps) {
	h := &questionHandler{
		log:     d.Logger,
		service: d.QuestionService,
	}

	g := handler.Group("questions")
	{
		g.GET("", h.getAll)
	}
}

func (h *questionHandler) getAll(c *gin.Context) {
	q, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - question - getAll - h.service.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, q)
}
