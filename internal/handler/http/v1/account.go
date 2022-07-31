package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
)

type accountHandler struct {
	cfg     *config.Session
	log     *zap.Logger
	v       ValidationModule
	account AccountService
}

func newAccountHandler(d *Deps) http.Handler {
	h := accountHandler{
		cfg:     &d.Config.Session,
		log:     d.Logger,
		v:       d.ValidationModule,
		account: d.AccountService,
	}

	r := chi.NewRouter()

	authenticator := mwAuthenticator(d.Logger, &d.Config.Session, d.SessionService)
	tokenRequired := mwTokenRequired(d.Logger, d.TokenService)

	r.Post("/", h.create)
	r.With(authenticator, tokenRequired).Delete("/", h.delete)

	r.Route("/verification", func(r chi.Router) {
		r.With(authenticator).Post("/", h.requestVerification)
		r.Put("/", h.verify)
	})

	r.Route("/password", func(r chi.Router) {
		r.Post("/", h.resetPassword)
		r.Put("/", h.setPassword)
	})

	return r
}

func (h *accountHandler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req account.CreateReq

	if err := h.v.ValidateRequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - create - ValidateRequestBody", zap.Error(err))
		writeDetailedError(w, http.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
		return
	}

	_, err := h.account.Create(r.Context(), req)
	if err != nil {
		h.log.Error("http - v1 - account - create - h.account.Create", zap.Error(err))

		switch {
		case errors.Is(err, account.ErrAlreadyExist):
			writeError(w, http.StatusConflict, account.ErrAlreadyExist)
			return
		case errors.Is(err, account.ErrForbiddenNickname):
			writeError(w, http.StatusBadRequest, account.ErrForbiddenNickname)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *accountHandler) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Error("http - v1 - account - delete - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = h.account.Delete(ctx, accountId); err != nil {
		h.log.Error("http - v1 - account - delete - h.account.Delete", zap.Error(err))

		if errors.Is(err, account.ErrNotDeleted) {
			writeError(w, http.StatusBadRequest, account.ErrAlreadyArchived)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

func (h *accountHandler) requestVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Error("http - v1 - account - requestVerification - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = h.account.RequestVerification(ctx, accountId); err != nil {
		h.log.Error("http - v1 - account - requestVerification - h.account.RequestVerification")

		if errors.Is(err, account.ErrAlreadyVerified) {
			writeError(w, http.StatusBadRequest, account.ErrAlreadyVerified)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *accountHandler) verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, account.ErrEmptyVerificationCode)
		return
	}

	if err := h.account.Verify(r.Context(), code); err != nil {
		h.log.Error("http - v1 - account - verify - h.account.Verify", zap.Error(err))

		if errors.Is(err, account.ErrAlreadyVerified) {
			writeError(w, http.StatusBadRequest, account.ErrAlreadyVerified)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *accountHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req account.ResetPasswordReq

	if err := h.v.ValidateRequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - resetPassword - ValidateRequestBody", zap.Error(err))
		writeDetailedError(w, http.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
		return
	}

	if err := h.account.ResetPassword(r.Context(), req.Login); err != nil {
		h.log.Error("http - v1 - account - resetPassword - h.account.ResetPassword", zap.Error(err))

		if errors.Is(err, account.ErrNotFound) {
			writeError(w, http.StatusNotFound, account.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *accountHandler) setPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query().Get("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, account.ErrEmptyPasswordToken)
		return
	}

	var req account.SetPasswordReq

	if err := h.v.ValidateRequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - setPassword - ValidateRequestBody", zap.Error(err))
		writeDetailedError(w, http.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
		return
	}

	if err := h.account.SetPassword(r.Context(), token, req.Password); err != nil {
		h.log.Error("http - v1 - account - setPassword - h.account.SetPassword", zap.Error(err))

		if errors.Is(err, account.ErrPasswordTokenExpired) ||
			errors.Is(err, account.ErrPasswordTokenNotFound) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
