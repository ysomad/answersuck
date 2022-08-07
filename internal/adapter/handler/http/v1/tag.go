package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/answersuck/host/internal/domain/tag"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type tagHandler struct {
	log      *zap.Logger
	validate validate
	service  tagService
}

func newTagMux(d *Deps) *chi.Mux {
	h := tagHandler{
		log:      d.Logger,
		validate: d.Validate,
		service:  d.TagService,
	}

	m := chi.NewMux()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	m.Post("/list", h.getAll)
	m.With(verificator).Post("/", h.createMultiple)

	return m
}

type tagGetAllReq struct {
	Filter struct {
		Name string `json:"name"`
	} `json:"filter"`
	LastId uint32 `json:"last_id"`
	Limit  uint64 `json:"limit"`
}

func (h *tagHandler) getAll(w http.ResponseWriter, r *http.Request) {
	var req tagGetAllReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Info("http - v1 - tag - getAll - RequestBody", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	tagList, err := h.service.GetAll(r.Context(), tag.NewListParams(req.LastId, req.Limit, tag.Filter{Name: req.Filter.Name}))
	if err != nil {
		h.log.Error("http - v1 - tag - getAll - h.service.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tagList)
}

type tagCreateMultipleReq struct {
	Tags []struct {
		Name       string `json:"name" binding:"required,gte=1,lte=32"`
		LanguageId uint8  `json:"language_id" binding:"required"`
	} `json:"tags" binding:"required,min=1,max=10"`
}

func (h *tagHandler) createMultiple(w http.ResponseWriter, r *http.Request) {
	var req tagCreateMultipleReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - tag - createMultiple - RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	t := make([]tag.Tag, 0, len(req.Tags))
	for _, rt := range req.Tags {
		t = append(t, tag.Tag{
			Name:       rt.Name,
			LanguageId: rt.LanguageId,
		})
	}

	tags, err := h.service.CreateMultiple(r.Context(), t)
	if err != nil {
		h.log.Error("http - v1 - tag - createMultiple - h.service.CreateMultiple", zap.Error(err))

		if errors.Is(err, tag.ErrLanguageIdNotFound) {
			writeErr(w, http.StatusBadRequest, tag.ErrLanguageIdNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, listResponse{tags})
}
