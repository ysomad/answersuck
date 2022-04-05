package v1

import (
	"errors"
	"fmt"
	"github.com/answersuck/vault/internal/service/repository"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/service"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/validation"
)

type accountHandler struct {
	t   validation.ErrorTranslator
	cfg *config.Aggregate
	log logging.Logger

	account  service.Account
	password service.AccountPassword
}

func newAccountHandler(handler *gin.RouterGroup, d *Deps) {
	h := &accountHandler{
		t:        d.ErrorTranslator,
		cfg:      d.Config,
		log:      d.Logger,
		account:  d.AccountService,
		password: d.AccountPasswordService,
	}

	accounts := handler.Group("accounts")
	{
		authenticated := accounts.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
		{
			authenticated.GET("", h.get)
			authenticated.DELETE("", tokenMiddleware(d.Logger, d.AuthService), h.archive)
			authenticated.POST("verify", h.requestVerification)
		}

		password := accounts.Group("password")
		{
			password.POST("forgot", h.passwordForgot)
			password.PATCH("reset", h.passwordReset)
		}

		accounts.PATCH("verify", h.verify)
		accounts.POST("", h.create)
	}
}

type accountCreateRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
	Password string `json:"password" binding:"required,gte=8,lte=64"`
}

func (h *accountHandler) create(c *gin.Context) {
	var r accountCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	_, err := h.account.Create(
		c.Request.Context(),
		&domain.Account{
			Email:    r.Email,
			Username: r.Username,
			Password: r.Password,
		},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - create: %w", err))

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

	c.Status(http.StatusCreated)
}

func (h *accountHandler) archive(c *gin.Context) {
	aid, err := accountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sid, err := sessionId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - sessionId: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = h.account.Delete(c.Request.Context(), aid, sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - h.account.Delete: %w", err))

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
	aid, err := accountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	acc, err := h.account.GetById(c.Request.Context(), aid)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - get - h.account.GetById: %w", err))

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
	aid, err := accountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - archive - accountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err = h.account.RequestVerification(c.Request.Context(), aid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - requestVerification - h.account.RequestVerification: %w", err))

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
	if !found || code == "" {
		abortWithError(c, http.StatusBadRequest, domain.ErrAccountEmptyVerificationCode, "")
		return
	}

	if err := h.account.Verify(c.Request.Context(), code, true); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - verify - h.account.Verify: %w", err))

		if errors.Is(err, repository.ErrNoAffectedRows) {
			abortWithError(c, http.StatusBadRequest, domain.ErrAccountAlreadyVerified, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

type accountPasswordForgotRequest struct {
	Login string `json:"login" binding:"required,email|alphanum"`
}

func (h *accountHandler) passwordForgot(c *gin.Context) {
	var r accountPasswordForgotRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	if err := h.password.RequestReset(c.Request.Context(), r.Login); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - passwordForgot - h.account.RequestPasswordReset: %w", err))

		if errors.Is(err, repository.ErrNotFound) {
			abortWithError(c, http.StatusNotFound, domain.ErrAccountNotFound, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

type accountPasswordResetRequest struct {
	Password string `json:"password" binding:"required,gte=8,lte=64"`
}

func (h *accountHandler) passwordReset(c *gin.Context) {
	var r accountPasswordResetRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		abortWithError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	t, found := c.GetQuery("token")
	if !found || t == "" {
		abortWithError(c, http.StatusBadRequest, domain.ErrAccountEmptyResetPasswordToken, "")
		return
	}

	if err := h.password.Reset(c.Request.Context(), t, r.Password); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - account - passwordReset - h.account.PasswordReset: %w", err))

		if errors.Is(err, domain.ErrAccountResetPasswordTokenExpired) || errors.Is(err, repository.ErrNotFound) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
