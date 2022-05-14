package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"

	"github.com/answersuck/vault/pkg/logging"
)

type AuthService interface {
	Login(ctx context.Context, login, password string, d domain.Device) (*domain.Session, error)
	Logout(ctx context.Context, sessionId string) error

	NewSecurityToken(ctx context.Context, accountId, password, audience string) (string, error)
	ParseSecurityToken(ctx context.Context, token, audience string) (string, error)
}

type authHandler struct {
	t       errorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	service AuthService
}

func newAuthHandler(r *gin.RouterGroup, d *Deps) {
	h := &authHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.AuthService,
	}

	auth := r.Group("auth")
	{

		auth.POST("login", deviceMiddleware(), h.login)
	}

	authenticated := auth.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
	{
		authenticated.POST("logout", h.logout)
		authenticated.POST("token", h.tokenCreate)
	}
}

func (h *authHandler) login(c *gin.Context) {
	var r dto.LoginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	d, err := getDevice(c)
	if err != nil {
		h.log.Error("http - v1 - auth - login - getDevice: %w", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	s, err := h.service.Login(c.Request.Context(), r.Login, r.Password, d)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - login - h.service.Login: %w", err))

		if errors.Is(err, domain.ErrAccountIncorrectPassword) || errors.Is(err, domain.ErrAccountNotFound) {
			abortWithError(c, http.StatusUnauthorized, domain.ErrAccountIncorrectCredentials, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.Session.CookieKey, s.Id, s.MaxAge, "", "", h.cfg.Cookie.Secure, h.cfg.Cookie.HTTPOnly)
	c.Status(http.StatusOK)
}

func (h *authHandler) logout(c *gin.Context) {
	sessionId := getSessionId(c)

	if err := h.service.Logout(c.Request.Context(), sessionId); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - logout - h.service.Logout: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.Session.CookieKey, "", -1, "", "", h.cfg.Cookie.Secure, h.cfg.Cookie.HTTPOnly)
	c.Status(http.StatusNoContent)
}

func (h *authHandler) tokenCreate(c *gin.Context) {
	var r dto.TokenCreateRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	t, err := h.service.NewSecurityToken(
		c.Request.Context(),
		accountId,
		r.Password,
		strings.ToLower(r.Audience),
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - h.service.NewToken: %w", err))

		if errors.Is(err, domain.ErrAccountIncorrectPassword) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dto.TokenCreateResponse{Token: t})
}
