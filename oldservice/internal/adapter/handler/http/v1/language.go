package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type languageHandler struct {
	log      *zap.Logger
	language languageService
}

func newLanguageHandler(d *Deps) *languageHandler {
	return &languageHandler{
		log:      d.Logger,
		language: d.LanguageService,
	}
}

func newLanguageMux(d *Deps) *chi.Mux {
	h := languageHandler{
		log:      d.Logger,
		language: d.LanguageService,
	}

	m := chi.NewMux()

	m.Get("/", h.getAll)

	return m
}

func (h *languageHandler) getAll(w http.ResponseWriter, r *http.Request) {
	langs, err := h.language.GetAll(r.Context())
	if err != nil {
		h.log.Info("http - v1 - language - getAll - h.language.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, listResponse{langs})
}
