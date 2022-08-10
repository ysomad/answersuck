package v1

import (
	"errors"
	"net/http"

	"github.com/answersuck/host/internal/domain/player"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type playerHandler struct {
	log      *zap.Logger
	validate validate
	player   playerService
}

func newPlayerMux(d *Deps) *chi.Mux {
	h := playerHandler{
		log:      d.Logger,
		validate: d.Validate,
		player:   d.PlayerService,
	}

	m := chi.NewMux()
	m.Get("/{nickname}", h.get)

	return m
}

func (h *playerHandler) get(w http.ResponseWriter, r *http.Request) {
	nickname := chi.URLParam(r, "nickname")
	if nickname == "" {
		writeErr(w, http.StatusNotFound, player.ErrNotFound)
		return
	}

	p, err := h.player.GetByNickname(r.Context(), nickname)
	if err != nil {
		h.log.Error("http - v1 - player - get - h.player.GetByNickname", zap.Error(err))

		if errors.Is(err, player.ErrNotFound) {
			writeErr(w, http.StatusNotFound, player.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}
