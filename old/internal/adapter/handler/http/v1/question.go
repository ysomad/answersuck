package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ysomad/answersuck-backend/internal/domain/question"
)

type questionHandler struct {
	log      *zap.Logger
	validate validate
	question questionService
}

func newQuestionMux(d *Deps) *chi.Mux {
	h := questionHandler{
		log:      d.Logger,
		validate: d.Validate,
		question: d.QuestionService,
	}

	m := chi.NewMux()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	m.With(verificator).Post("/", h.create)
	m.Get("/{question_id}", h.get)
	m.Post("/list", h.getAll)

	return m
}

type questionCreateReq struct {
	Text       string `json:"text" validate:"required,gte=1,lte=200"`
	AnswerId   uint32 `json:"answer_id" validate:"required"`
	MediaId    string `json:"media_id" validate:"omitempty,uuid4"`
	LanguageId uint8  `json:"language_id" validate:"required"`
}

type questionCreateRes struct {
	Id uint32 `json:"id"`
}

func (h *questionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req questionCreateReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - question - create - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - question - create - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	questionId, err := h.question.Create(r.Context(), question.CreateDTO{
		Text:       req.Text,
		AnswerId:   req.AnswerId,
		AccountId:  accountId,
		LanguageId: req.LanguageId,
	})
	if err != nil {
		h.log.Error("http - v1 - question - create - h.question.Create", zap.Error(err))

		if errors.Is(err, question.ErrForeignKeyViolation) {
			writeErr(w, http.StatusBadRequest, question.ErrForeignKeyViolation)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, questionCreateRes{questionId})
}

func (h *questionHandler) get(w http.ResponseWriter, r *http.Request) {
	s := chi.URLParam(r, "question_id")
	if s == "" {
		writeErr(w, http.StatusNotFound, question.ErrNotFound)
		return
	}

	questionId, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		h.log.Info("http - v1 - question - get - strconv.ParseUint", zap.Error(err))
		writeErr(w, http.StatusNotFound, question.ErrNotFound)
		return
	}

	q, err := h.question.GetById(r.Context(), uint32(questionId))
	if err != nil {
		h.log.Error("http - v1 - question - get - h.question.GetById", zap.Error(err))

		if errors.Is(err, question.ErrNotFound) {
			writeErr(w, http.StatusNotFound, question.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, q)
}

type questionGetAllReq struct {
	Filter struct {
		Text       string `json:"text"`
		Author     string `json:"author"`
		LanguageId uint   `json:"language_id"`
	} `json:"filter"`
	LastId uint32 `json:"last_id"`
	Limit  uint64 `json:"limit"`
}

func (h *questionHandler) getAll(w http.ResponseWriter, r *http.Request) {
	var req questionGetAllReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Info("http - v1 - answer - getAll - Decode", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	questionList, err := h.question.GetAll(r.Context(), question.NewListParams(req.LastId, req.Limit, question.Filter{
		Text:       req.Filter.Text,
		Author:     req.Filter.Author,
		LanguageId: req.Filter.LanguageId,
	}))
	if err != nil {
		h.log.Error("http - v1 - question - getAll - h.question.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, questionList)
}
