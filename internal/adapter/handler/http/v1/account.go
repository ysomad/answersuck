package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/config"
	"github.com/answersuck/host/internal/domain/account"
)

type accountHandler struct {
	cfg      *config.Session
	log      *zap.Logger
	validate validate
	account  accountService
}

func newAccountMux(d *Deps) *chi.Mux {
	h := accountHandler{
		cfg:      &d.Config.Session,
		log:      d.Logger,
		validate: d.Validate,
		account:  d.AccountService,
	}

	m := chi.NewMux()

	authenticator := mwAuthenticator(d.Logger, &d.Config.Session, d.SessionService)
	tokenRequired := mwTokenRequired(d.Logger, d.TokenService)

	m.Post("/", h.create)
	m.With(authenticator, tokenRequired).Delete("/", h.delete)

	m.Route("/verification", func(r chi.Router) {
		r.With(authenticator).Post("/", h.requestVerification)
		r.Put("/", h.verify)
	})

	m.Route("/password", func(r chi.Router) {
		r.With(authenticator).Patch("/", h.updatePassword)
		r.Post("/", h.resetPassword)
		r.Put("/", h.setPassword)
	})

	return m
}

type accountCreateReq struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Nickname string `json:"nickname" validate:"required,alphanum,gte=4,lte=25"`
	Password string `json:"password" validate:"required,gte=10,lte=128"`
}

func (h *accountHandler) create(w http.ResponseWriter, r *http.Request) {
	var req accountCreateReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - create - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	_, err := h.account.Create(r.Context(), req.Email, req.Nickname, req.Password)
	if err != nil {
		h.log.Error("http - v1 - account - create - h.account.Create", zap.Error(err))

		switch {
		case errors.Is(err, account.ErrAlreadyExist):
			writeErr(w, http.StatusConflict, account.ErrAlreadyExist)
			return
		case errors.Is(err, account.ErrForbiddenNickname):
			writeErr(w, http.StatusBadRequest, account.ErrForbiddenNickname)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *accountHandler) delete(w http.ResponseWriter, r *http.Request) {
	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - account - delete - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = h.account.Delete(r.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - delete - h.account.Delete", zap.Error(err))

		if errors.Is(err, account.ErrNotDeleted) {
			writeErr(w, http.StatusBadRequest, account.ErrAlreadyArchived)
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
	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - account - requestVerification - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = h.account.RequestVerification(r.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - requestVerification - h.account.RequestVerification")

		if errors.Is(err, account.ErrAlreadyVerified) {
			writeErr(w, http.StatusBadRequest, account.ErrAlreadyVerified)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *accountHandler) verify(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeErr(w, http.StatusBadRequest, account.ErrEmptyVerificationCode)
		return
	}

	if err := h.account.Verify(r.Context(), code); err != nil {
		h.log.Error("http - v1 - account - verify - h.account.Verify", zap.Error(err))

		if errors.Is(err, account.ErrAlreadyVerified) {
			writeErr(w, http.StatusBadRequest, account.ErrAlreadyVerified)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type accountUpdatePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,gte=10,lte=128,nefield=OldPassword"`
}

func (h *accountHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Error("http - v1 - account - updatePassword - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req accountUpdatePasswordReq
	if err = h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - updatePassword - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	if err = h.account.UpdatePassword(ctx, accountId, req.OldPassword, req.NewPassword); err != nil {
		h.log.Error("http - v1 - account - updatePassword - h.account.UpdatePassword", zap.Error(err))

		if errors.Is(err, account.ErrInvalidPassword) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type accountResetPasswordReq struct {
	Login string `json:"login" validate:"required,email|alphanum"`
}

func (h *accountHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req accountResetPasswordReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - resetPassword - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	if err := h.account.ResetPassword(r.Context(), req.Login); err != nil {
		h.log.Error("http - v1 - account - resetPassword - h.account.ResetPassword", zap.Error(err))

		if errors.Is(err, account.ErrNotFound) {
			writeErr(w, http.StatusNotFound, account.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type accountSetPasswordReq struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,gte=10,lte=128"`
}

func (h *accountHandler) setPassword(w http.ResponseWriter, r *http.Request) {
	var req accountSetPasswordReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - account - setPassword - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	if len(req.Token) != account.PasswordTokenLen {
		h.log.Info("http - v1 - account - setPassword", zap.Error(account.ErrPasswordTokenInvalid))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := h.account.SetPassword(r.Context(), req.Token, req.Password); err != nil {
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
