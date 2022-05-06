package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
	repository "github.com/answersuck/vault/internal/repository/psql"

	"github.com/answersuck/vault/pkg/logging"
)

type accountService interface {
	Create(ctx context.Context, req dto.AccountCreateRequest) (*domain.Account, error)
	GetById(ctx context.Context, aid string) (*domain.Account, error)
	Delete(ctx context.Context, aid, sid string) error
	RequestVerification(ctx context.Context, aid string) error
	Verify(ctx context.Context, code string, verified bool) error
	RequestPasswordReset(ctx context.Context, login string) error
	PasswordReset(ctx context.Context, token, password string) error
}

type accountHandler struct {
	t   errorTranslator
	cfg *config.Aggregate
	log logging.Logger

	service accountService
}

func newAccountHandler(handler *gin.RouterGroup, d *Deps) {
	h := &accountHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.AccountService,
	}

	accounts := handler.Group("accounts")
	{
		authenticated := accounts.Group("",
			sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService),
			accountParamMiddleware(d.Logger))
		{
			authenticated.GET(":accountId", h.get)
			authenticated.DELETE(":accountId", tokenMiddleware(d.Logger, d.AuthService), h.archive)
		}

		passwordReset := accounts.Group("password/reset")
		{
			passwordReset.POST("", h.passwordForgot)
			passwordReset.PUT("", h.passwordReset)
		}

		accounts.POST(":accountId/verification",
			sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService),
			accountParamMiddleware(d.Logger),
			h.requestVerification)
		accounts.PUT("verification", h.verify)

		accounts.POST("", h.create)
	}
}

func (h *accountHandler) create(c *gin.Context) {
	var r dto.AccountCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	a, err := h.service.Create(c.Request.Context(), r)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - create - h.service.Create: %w", err))

		if errors.Is(err, repository.ErrUniqueViolation) {
			abortWithError(c, http.StatusConflict, domain.ErrAccountAlreadyExist, "")
			return
		}

		if errors.Is(err, domain.ErrAccountForbiddenUsername) {
			abortWithError(c, http.StatusBadRequest, domain.ErrAccountForbiddenUsername, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *accountHandler) archive(c *gin.Context) {
	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.service.Delete(c.Request.Context(), aid, getSessionId(c)); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - h.service.Delete: %w", err))

		if errors.Is(err, repository.ErrNotFound) || errors.Is(err, repository.ErrNoAffectedRows) {
			abortWithError(c, http.StatusNotFound, domain.ErrAccountNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.Session.CookieKey, "", -1, "", "", h.cfg.Cookie.Secure, h.cfg.Cookie.HTTPOnly)
	c.Status(http.StatusNoContent)
}

func (h *accountHandler) get(c *gin.Context) {
	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - getAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	acc, err := h.service.GetById(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - get - h.service.GetById: %w", err))

		if errors.Is(err, repository.ErrNotFound) {
			abortWithError(c, http.StatusNotFound, domain.ErrAccountNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (h *accountHandler) requestVerification(c *gin.Context) {
	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.service.RequestVerification(c.Request.Context(), aid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - requestVerification - h.service.RequestVerification: %w", err))

		if errors.Is(err, repository.ErrNotFound) {
			abortWithError(c, http.StatusNotFound, domain.ErrAccountNotFound, "")
			return
		}

		if errors.Is(err, domain.ErrAccountAlreadyVerified) {
			abortWithError(c, http.StatusBadRequest, domain.ErrAccountAlreadyVerified, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *accountHandler) verify(c *gin.Context) {
	code, found := c.GetQuery("code")
	if !found {
		abortWithError(c, http.StatusBadRequest, domain.ErrAccountEmptyVerificationCode, "")
		return
	}

	if err := h.service.Verify(c.Request.Context(), code, true); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - verify - h.service.Verify: %w", err))

		if errors.Is(err, repository.ErrNoAffectedRows) {
			abortWithError(c, http.StatusBadRequest, domain.ErrAccountAlreadyVerified, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *accountHandler) passwordForgot(c *gin.Context) {
	var r dto.AccountPasswordForgotRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	if err := h.service.RequestPasswordReset(c.Request.Context(), r.Login); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - passwordForgot - h.service.RequestPasswordReset: %w", err))

		if errors.Is(err, repository.ErrNotFound) {
			abortWithError(c, http.StatusNotFound, domain.ErrAccountNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *accountHandler) passwordReset(c *gin.Context) {
	var r dto.AccountPasswordResetRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, found := c.GetQuery("token")
	if !found || t == "" {
		abortWithError(c, http.StatusBadRequest, domain.ErrAccountEmptyResetPasswordToken, "")
		return
	}

	if err := h.service.PasswordReset(c.Request.Context(), t, r.Password); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - passwordReset - h.service.PasswordReset: %w", err))

		if errors.Is(err, domain.ErrAccountResetPasswordTokenExpired) || errors.Is(err, repository.ErrNotFound) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
