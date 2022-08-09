package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/answersuck/host/internal/domain/answer"
	"github.com/answersuck/host/internal/domain/media"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type answerHandler struct {
	log      *zap.Logger
	validate validate
	answer   answerService
}

func newAnswerMux(d *Deps) *chi.Mux {
	h := answerHandler{
		log:      d.Logger,
		validate: d.Validate,
		answer:   d.AnswerService,
	}

	m := chi.NewMux()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	m.Post("/list", h.getAll)
	m.With(verificator).Post("/", h.create)

	return m
}

type answerCreateReq struct {
	Text       string `json:"text" validate:"required,gte=1,lte=100"`
	MediaId    string `json:"media_id" validate:"omitempty,uuid4"`
	LanguageId uint   `json:"language_id" validate:"required"`
}

func (h *answerHandler) create(w http.ResponseWriter, r *http.Request) {
	var req answerCreateReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - answer - create - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	a, err := h.answer.Create(r.Context(), answer.Answer{
		Text:       req.Text,
		MediaId:    &req.MediaId,
		LanguageId: uint(req.LanguageId),
	})
	if err != nil {
		h.log.Error("http - v1 - answer - create - h.answer.Create", zap.Error(err))

		switch {
		case errors.Is(err, answer.ErrMediaTypeNotAllowed):
			writeErr(w, http.StatusBadRequest, answer.ErrMediaTypeNotAllowed)
			return
		case errors.Is(err, answer.ErrLanguageNotFound):
			writeErr(w, http.StatusBadRequest, answer.ErrLanguageNotFound)
			return
		case errors.Is(err, media.ErrNotFound):
			writeErr(w, http.StatusBadRequest, media.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, a)
}

type answerGetAllReq struct {
	Filter struct {
		Text       string `json:"text"`
		LanguageId uint   `json:"language_id"`
	} `json:"filter"`
	LastId uint32 `json:"last_id"`
	Limit  uint64 `json:"limit"`
}

func (h *answerHandler) getAll(w http.ResponseWriter, r *http.Request) {
	var req answerGetAllReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Info("http - v1 - answer - getAll - Decode", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	answers, err := h.answer.GetAll(
		r.Context(),
		answer.NewListParams(req.LastId, req.Limit, answer.Filter{
			Text:       req.Filter.Text,
			LanguageId: req.Filter.LanguageId,
		},
		))
	if err != nil {
		h.log.Info("http - v1 - answer - getAll - h.answer.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, answers)
}
