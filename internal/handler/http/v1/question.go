package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/service"
	"github.com/answersuck/vault/pkg/logging"
)

type questionHandler struct {
	log      logging.Logger
	question service.Question
}

func newQuestionHandler(handler *gin.RouterGroup, d *Deps) {
	h := &questionHandler{
		log:      d.Logger,
		question: d.QuestionService,
	}

	g := handler.Group("questions")
	{
		g.GET("", h.getAll)
	}
}

func (h *questionHandler) getAll(c *gin.Context) {
	q, err := h.question.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - question - getAll - h.question.GetAll: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, q)
}
