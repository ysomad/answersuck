package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"

	"github.com/answersuck/vault/pkg/logging"
)

type AccountService interface {
	Create(ctx context.Context, req account.CreateRequest) (*account.Account, error)
	GetById(ctx context.Context, accountId string) (*account.Account, error)
	Delete(ctx context.Context, accountId string) error
	RequestVerification(ctx context.Context, accountId string) error
	Verify(ctx context.Context, code string, verified bool) error
	ResetPassword(ctx context.Context, login string) error
	SetPassword(ctx context.Context, token, password string) error
}

type accountHandler struct {
	t   ErrorTranslator
	cfg *config.Aggregate
	log logging.Logger

	service AccountService
}

func newAccountHandler(r *gin.RouterGroup, d *Deps) {
	h := &accountHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.AccountService,
	}

	accounts := r.Group("accounts")
	{
		accounts.POST("", h.create)
	}

	authenticated := accounts.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
	{
		authenticated.GET("", h.get)
		authenticated.DELETE("", tokenMiddleware(d.Logger, d.AuthService), h.delete)
	}

	verification := accounts.Group("verification")
	{
		verification.POST("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService), h.requestVerification)
		verification.PUT("", h.verify)
	}

	password := accounts.Group("password")
	{
		password.POST("", h.resetPassword)
		password.PUT("", h.setPassword)
	}
}

func (h *accountHandler) create(c *gin.Context) {
	var r account.CreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	a, err := h.service.Create(c.Request.Context(), r)
	if err != nil {
		h.log.Error("http - v1 - account - create - h.service.Create: %w", err)

		switch {
		case errors.Is(err, account.ErrAlreadyExist):
			abortWithError(c, http.StatusConflict, account.ErrAlreadyExist, "")
			return
		case errors.Is(err, account.ErrForbiddenUsername):
			abortWithError(c, http.StatusBadRequest, account.ErrForbiddenUsername, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *accountHandler) delete(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - account - delete - getAccountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.service.Delete(c.Request.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - delete - h.service.Delete: %w", err)

		if errors.Is(err, account.ErrNotArchived) {
			abortWithError(c, http.StatusBadRequest, account.ErrAlreadyArchived, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.Session.CookieKey, "", -1, "", "", h.cfg.Cookie.Secure, h.cfg.Cookie.HTTPOnly)
	c.Status(http.StatusNoContent)
}

func (h *accountHandler) get(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - account - get - getAccountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	a, err := h.service.GetById(c.Request.Context(), accountId)
	if err != nil {
		h.log.Error("http - v1 - account - get - h.service.GetById: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *accountHandler) requestVerification(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - account - requestVerification - getAccountId: %w", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.service.RequestVerification(c.Request.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - requestVerification - h.service.RequestVerification: %w", err)

		if errors.Is(err, account.ErrAlreadyVerified) {
			abortWithError(c, http.StatusBadRequest, account.ErrAlreadyVerified, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

const (
	queryCode = "code"
)

func (h *accountHandler) verify(c *gin.Context) {
	code, found := c.GetQuery(queryCode)
	if !found {
		abortWithError(c, http.StatusBadRequest, account.ErrEmptyVerificationCode, "")
		return
	}

	if err := h.service.Verify(c.Request.Context(), code, true); err != nil {
		h.log.Error("http - v1 - account - verify - h.service.Verify: %w", err)

		if errors.Is(err, account.ErrAlreadyVerified) {
			abortWithError(c, http.StatusBadRequest, account.ErrAlreadyVerified, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *accountHandler) resetPassword(c *gin.Context) {
	var r account.ResetPasswordRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), r.Login); err != nil {
		h.log.Error("http - v1 - account - resetPassword - h.service.ResetPassword: %w", err)

		if errors.Is(err, account.ErrNotFound) {
			abortWithError(c, http.StatusNotFound, account.ErrNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

const (
	queryToken = "token"
)

func (h *accountHandler) setPassword(c *gin.Context) {
	var r account.SetPasswordRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, found := c.GetQuery(queryToken)
	if !found || t == "" {
		abortWithError(c, http.StatusBadRequest, account.ErrEmptyPasswordResetToken, "")
		return
	}

	if err := h.service.SetPassword(c.Request.Context(), t, r.Password); err != nil {
		h.log.Error("http - v1 - account - setPassword - h.service.SetPassword: %w", err)

		if errors.Is(err, account.ErrPasswordResetTokenExpired) ||
			errors.Is(err, account.ErrPasswordResetTokenNotFound) {

			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
