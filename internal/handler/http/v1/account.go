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

	"github.com/answersuck/vault/pkg/logging"
)

type AccountService interface {
	Create(ctx context.Context, req dto.AccountCreateRequest) (*domain.Account, error)
	GetById(ctx context.Context, accountId string) (*domain.Account, error)
	Delete(ctx context.Context, accountId string) error
	RequestVerification(ctx context.Context, accountId string) error
	Verify(ctx context.Context, code string, verified bool) error
	PasswordReset(ctx context.Context, login string) error
	PasswordSet(ctx context.Context, token, password string) error
}

type accountHandler struct {
	t   errorTranslator
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
		authenticated.DELETE("", tokenMiddleware(d.Logger, d.AuthService), h.archive)
	}

	verification := accounts.Group("verification")
	{
		verification.POST("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService), h.requestVerification)
		verification.PUT("", h.verify)
	}

	password := accounts.Group("password")
	{
		password.POST("", h.passwordReset)
		password.PUT("", h.passwordSet)
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

		if errors.Is(err, domain.ErrAccountAlreadyExist) {
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

	c.JSON(http.StatusOK, dto.AccountCreateResponse{
		Id:        a.Id,
		Email:     a.Email,
		Username:  a.Username,
		Verified:  a.Verified,
		Archived:  a.Archived,
		AvatarURL: a.AvatarURL,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	})
}

func (h *accountHandler) archive(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.service.Delete(c.Request.Context(), accountId); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - h.service.Delete: %w", err))

		if errors.Is(err, domain.ErrAccountNotArchived) {
			abortWithError(c, http.StatusBadRequest, domain.ErrAccountAlreadyArchived, "")
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
		h.log.Error(fmt.Errorf("http - v1 - account - archive - getAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	a, err := h.service.GetById(c.Request.Context(), accountId)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - get - h.service.GetById: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dto.AccountCreateResponse{
		Id:        a.Id,
		Email:     a.Email,
		Username:  a.Username,
		Verified:  a.Verified,
		Archived:  a.Archived,
		AvatarURL: a.AvatarURL,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	})
}

func (h *accountHandler) requestVerification(c *gin.Context) {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.service.RequestVerification(c.Request.Context(), accountId); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - requestVerification - h.service.RequestVerification: %w", err))

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

		if errors.Is(err, domain.ErrAccountAlreadyVerified) {
			abortWithError(c, http.StatusBadRequest, domain.ErrAccountAlreadyVerified, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *accountHandler) passwordReset(c *gin.Context) {
	var r dto.AccountPasswordResetRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	if err := h.service.PasswordReset(c.Request.Context(), r.Login); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - passwordForgot - h.service.RequestPasswordReset: %w", err))

		if errors.Is(err, domain.ErrAccountNotFound) {
			abortWithError(c, http.StatusNotFound, domain.ErrAccountNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *accountHandler) passwordSet(c *gin.Context) {
	var r dto.AccountPasswordSetRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, found := c.GetQuery("token")
	if !found || t == "" {
		abortWithError(c, http.StatusBadRequest, domain.ErrAccountEmptyResetPasswordToken, "")
		return
	}

	if err := h.service.PasswordSet(c.Request.Context(), t, r.Password); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - passwordReset - h.service.PasswordReset: %w", err))

		if errors.Is(err, domain.ErrAccountPasswordResetTokenExpired) ||
			errors.Is(err, domain.ErrAccountPasswordResetTokenNotFound) {

			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
