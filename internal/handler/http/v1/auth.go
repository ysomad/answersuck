package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/auth"
)

type authHandler struct {
	cfg     *config.Session
	log     *zap.Logger
	v       ValidationModule
	service LoginService
	token   TokenService
}

func newAuthHandler(d *Deps) http.Handler {
	h := authHandler{
		cfg:     &d.Config.Session,
		log:     d.Logger,
		v:       d.ValidationModule,
		service: d.LoginService,
		token:   d.TokenService,
	}

	r := chi.NewRouter()

	authenticator := mwAuthenticator(d.Logger, &d.Config.Session, d.SessionService)

	r.With(mwDeviceCtx).Post("/login", h.login)
	r.With(authenticator).Post("/token", h.createToken)
	r.Post("/logout", h.logout)

	return r
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := r.Cookie(h.cfg.CookieName)
	if err == nil {
		writeError(w, http.StatusBadRequest, auth.ErrAlreadyLoggedIn)
		return
	}

	var req auth.LoginReq

	if err = h.v.ValidateRequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - auth - login - h.v.ValidateRequestBody", zap.Error(err))
		writeDetailedError(w, http.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
		return
	}

	ctx := r.Context()

	d, err := getDevice(ctx)
	if err != nil {
		h.log.Error("http - v1 - auth - login - getDevice", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s, err := h.service.Login(ctx, req.Login, req.Password, d)
	if err != nil {
		h.log.Error("http - v1 - auth - login - h.service.Login", zap.Error(err))

		if errors.Is(err, auth.ErrIncorrectAccountPassword) ||
			errors.Is(err, account.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, auth.ErrIncorrectCredentials)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.CookieName,
		Value:    s.Id,
		Path:     h.cfg.CookiePath,
		MaxAge:   s.MaxAge,
		Secure:   h.cfg.CookieSecure,
		HttpOnly: h.cfg.CookieHTTPOnly,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *authHandler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.CookieName,
		Value:    "",
		Path:     h.cfg.CookiePath,
		MaxAge:   -1,
		Secure:   h.cfg.CookieSecure,
		HttpOnly: h.cfg.CookieHTTPOnly,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *authHandler) createToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req auth.TokenCreateReq

	if err := h.v.ValidateRequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - auth - createToken - h.v.ValidateRequestBody", zap.Error(err))
		writeDetailedError(w, http.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
		return
	}

	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Error("http - v1 - auth - createToken - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	t, err := h.token.Create(ctx, accountId, req.Password)
	if err != nil {
		h.log.Error("http - v1 - auth - createToken - h.token.Create", zap.Error(err))

		if errors.Is(err, auth.ErrIncorrectAccountPassword) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, auth.TokenCreateResp{Token: t})
	w.WriteHeader(http.StatusOK)
}
