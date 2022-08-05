package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/session"
)

type sessionHandler struct {
	log      *zap.Logger
	validate validate
	session  sessionService
}

func newSessionMux(d *Deps) *chi.Mux {
	h := sessionHandler{
		log:      d.Logger,
		validate: d.Validate,
		session:  d.SessionService,
	}

	m := chi.NewMux()

	authenticator := mwAuthenticator(d.Logger, &d.Config.Session, d.SessionService)
	tokenRequired := mwTokenRequired(d.Logger, d.TokenService)

	m.Route("/", func(r chi.Router) {
		r.Use(authenticator)
		r.Get("/", h.getAll)
		r.With(tokenRequired).Delete("/", h.terminateAll)
	})

	m.With(authenticator, tokenRequired).Delete("/{sessionId}", h.terminate)

	return m
}

func (h *sessionHandler) getAll(w http.ResponseWriter, r *http.Request) {
	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - session - getAll - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := h.session.GetAll(r.Context(), accountId)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - h.service.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, listResponse{s})
}

func (h *sessionHandler) terminate(w http.ResponseWriter, r *http.Request) {
	currSessionId, err := getSessionId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - session - terminate - getSessionId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if currSessionId == sessionId {
		writeErr(w, http.StatusBadRequest, session.ErrCannotBeTerminated)
		return
	}

	if err := h.session.Terminate(r.Context(), sessionId); err != nil {
		h.log.Error("http - v1 - session - terminate - h.service.Terminate", zap.Error(err))

		if errors.Is(err, session.ErrNotFound) {
			writeErr(w, http.StatusNotFound, session.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(w http.ResponseWriter, r *http.Request) {
	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currSessionId, err := getSessionId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - getSessionId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = h.session.TerminateAllWithExcept(r.Context(), accountId, currSessionId); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.service.TerminateWithExcept", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
