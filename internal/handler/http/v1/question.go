package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
	repository "github.com/answersuck/vault/internal/repository/psql"

	"github.com/answersuck/vault/pkg/logging"
)

type questionService interface {
	Create(ctx context.Context, dto *dto.QuestionCreate) (*domain.Question, error)
	GetAll(ctx context.Context) ([]*domain.Question, error)
}

type questionHandler struct {
	t       errorTranslator
	log     logging.Logger
	service questionService
}

func newQuestionHandler(handler *gin.RouterGroup, d *Deps) {
	h := &questionHandler{
		t:       d.ErrorTranslator,
		log:     d.Logger,
		service: d.QuestionService,
	}

	g := handler.Group("questions")
	{
		g.GET("", h.getAll)
	}
}

func (h *questionHandler) create(c *gin.Context) {
	var r dto.QuestionCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - question - create - getAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	q, err := h.service.Create(c.Request.Context(), &dto.QuestionCreate{
		Question:      r.Question,
		MediaId:       r.MediaId,
		Answer:        r.Answer,
		AnswerImageId: r.AnswerImageId,
		LanguageId:    r.LanguageId,
		AccountId:     aid,
	})
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - question - create - h.service.Create :%w", err))

		if errors.Is(err, repository.ErrForeignKeyViolation) {
			abortWithError(c, http.StatusBadRequest, repository.ErrForeignKeyViolation, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, q)
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
