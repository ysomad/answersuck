package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/domain/question"

	"github.com/answersuck/vault/pkg/logging"
)

type QuestionService interface {
	Create(ctx context.Context, q *question.Question) (*question.Question, error)
	GetById(ctx context.Context, questionId int) (*question.Detailed, error)
	GetAll(ctx context.Context) ([]question.Minimized, error)
}

type questionHandler struct {
	t       ErrorTranslator
	log     logging.Logger
	service QuestionService
}

func newQuestionHandler(r *gin.RouterGroup, d *Deps) {
	h := &questionHandler{
		t:       d.ErrorTranslator,
		log:     d.Logger,
		service: d.QuestionService,
	}

	questions := r.Group("questions")
	{
		questions.GET("", h.getAll)
		questions.GET(":questionId", h.getById)
	}

	authenticated := questions.Group("",
		sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService),
		protectionMiddleware(d.Logger))
	{
		authenticated.POST("", h.create)
	}
}

func (h *questionHandler) create(c *gin.Context) {
	var r question.CreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - question - create - getAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	q := question.Question{
		Text:       r.Text,
		AnswerId:   r.AnswerId,
		AccountId:  accountId,
		LanguageId: r.LanguageId,
	}

	if r.MediaId != "" {
		q.MediaId = &r.MediaId
	}

	res, err := h.service.Create(c.Request.Context(), &q)
	if err != nil {
		h.log.Error("http - v1 - question - create - h.service.Create :%w", err)

		if errors.Is(err, question.ErrForeignKeyViolation) {
			abortWithError(c, http.StatusBadRequest, question.ErrForeignKeyViolation, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *questionHandler) getAll(c *gin.Context) {
	qs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error("http - v1 - question - getAll - h.service.GetAll: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, qs)
}

func (h *questionHandler) getById(c *gin.Context) {
	s := c.Param("questionId")

	questionId, err := strconv.Atoi(s)
	if err != nil {
		h.log.Error("http - v1 - question - getById: %w", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	q, err := h.service.GetById(c.Request.Context(), questionId)
	if err != nil {
		h.log.Error("http - v1 - question - getById: %w", err)

		if errors.Is(err, question.ErrNotFound) {
			abortWithError(c, http.StatusNotFound, question.ErrNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, q)
}
