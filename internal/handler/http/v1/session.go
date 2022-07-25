package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/domain/session"
)

type sessionHandler struct {
	log     *zap.Logger
	v       ValidationModule
	session SessionService
}

func newSessionHandler(d *Deps) http.Handler {
	h := sessionHandler{
		log:     d.Logger,
		v:       d.ValidationModule,
		session: d.SessionService,
	}

	r := chi.NewRouter()

	authenticator := mwAuthenticator(d.Logger, &d.Config.Session, d.SessionService)
	tokenRequired := mwTokenRequired(d.Logger, d.TokenService)

	r.Route("/", func(r chi.Router) {
		r.Use(authenticator)
		r.Get("/", h.getAll)
		r.With(tokenRequired).Delete("/", h.terminateAll)
	})

	r.With(authenticator, tokenRequired).Delete("/{sessionId}", h.terminate)

	return r
}

func (h *sessionHandler) getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := h.session.GetAll(ctx, accountId)
	if err != nil {
		h.log.Error("http - v1 - session - getAll - h.service.GetAll", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeList(w, s)
	w.WriteHeader(http.StatusOK)
}

func (h *sessionHandler) terminate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	currSessionId, err := getSessionId(ctx)
	if err != nil {
		h.log.Error("http - v1 - session - terminate - getSessionId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId := chi.URLParam(r, "sessionId")
	if currSessionId == sessionId {
		writeError(w, http.StatusBadRequest, session.ErrCannotBeTerminated)
		return
	}

	if err := h.session.Terminate(ctx, sessionId); err != nil {
		h.log.Error("http - v1 - session - terminate - h.service.Terminate", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *sessionHandler) terminateAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId, err := getSessionId(ctx)
	if err != nil {
		h.log.Error("http - v1 - session - terminateAll - getSessionId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = h.session.TerminateWithExcept(ctx, accountId, sessionId); err != nil {
		h.log.Error("http - v1 - session - terminateAll - h.service.TerminateWithExcept", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
