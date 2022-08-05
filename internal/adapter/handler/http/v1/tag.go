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

	m.Get("/", h.getAll)
	m.With(verificator).Post("/", h.createMultiple)

	return m
}

type tagCreateMultipleReq struct {
	Tags []struct {
		Name       string `json:"name" binding:"required,gte=1,lte=32"`
		LanguageId uint32 `json:"language_id" binding:"required"`
	} `json:"tags" binding:"required,min=1,max=10"`
}

func (r *tagCreateMultipleReq) toTagList() ([]tag.Tag, error) {
	b, err := json.Marshal(r.Tags)
	if err != nil {
		return nil, err
	}

	var tags []tag.Tag
	if err = json.Unmarshal(b, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (h *tagHandler) createMultiple(w http.ResponseWriter, r *http.Request) {
	var req tagCreateMultipleReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - tag - createMultiple - RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	t, err := req.toTagList()
	if err != nil {
		h.log.Error("http - v1 - createMultiple - req.toTagList", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Error("TEST", zap.Any("tags", t))

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

func (h *tagHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.GetAll(r.Context())
	if err != nil {
		h.log.Error("http - v1 - tag - getAll - h.service.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, listResponse{tags})
}
