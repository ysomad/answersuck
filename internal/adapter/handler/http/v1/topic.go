package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/answersuck/host/internal/domain/topic"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type topicHandler struct {
	log      *zap.Logger
	validate validate
	topic    topicService
}

func newTopicMux(d *Deps) *chi.Mux {
	h := topicHandler{
		log:      d.Logger,
		validate: d.Validate,
		topic:    d.TopicService,
	}

	m := chi.NewMux()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	m.With(verificator).Post("/", h.create)
	m.Post("/list", h.getAll)

	return m
}

type topicCreateReq struct {
	Name       string `json:"name" binding:"required,gte=4,lte=50"`
	LanguageId uint   `json:"language_id" binding:"required"`
}

func (h *topicHandler) create(w http.ResponseWriter, r *http.Request) {
	var req topicCreateReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - topic - create - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	t, err := h.topic.Create(r.Context(), topic.Topic{
		Name:       req.Name,
		LanguageId: req.LanguageId,
	})
	if err != nil {
		h.log.Error("http - v1 - topic - create - h.topic.Create", zap.Error(err))

		if errors.Is(err, topic.ErrLanguageNotFound) {
			writeErr(w, http.StatusBadRequest, topic.ErrLanguageNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, t)
}

type topicGetAllReq struct {
	Filter struct {
		Name       string `json:"name"`
		LanguageId uint   `json:"language_id"`
	} `json:"filter"`
	LastId uint32 `json:"last_id"`
	Limit  uint64 `json:"limit"`
}

func (h *topicHandler) getAll(w http.ResponseWriter, r *http.Request) {
	var req topicGetAllReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Info("http - v1 - topic - getAll - Decode", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	topicList, err := h.topic.GetAll(r.Context(), topic.NewListParams(req.LastId, req.Limit, topic.Filter{
		Name:       req.Filter.Name,
		LanguageId: req.Filter.LanguageId,
	}))
	if err != nil {
		h.log.Error("http - v1 - topic - getAll - h.topic.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, topicList)
}
