package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/app/auth"
	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

type AuthService interface {
	Login(ctx context.Context, login, password string, d session.Device) (*session.Session, error)
	Logout(ctx context.Context, sessionId string) error

	NewToken(ctx context.Context, accountId, password, audience string) (string, error)
	ParseToken(ctx context.Context, token, audience string) (string, error)
}

type authHandler struct {
	t       ErrorTranslator
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
		authenticated.POST("token", h.createToken)
	}
}

func (h *authHandler) login(c *gin.Context) {
	if _, err := c.Cookie(h.cfg.Session.CookieKey); !errors.Is(err, http.ErrNoCookie) {
		abortWithError(c, http.StatusBadRequest, auth.ErrAlreadyLoggedIn, "")
		return
	}

	var r auth.LoginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
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

		if errors.Is(err, account.ErrIncorrectPassword) || errors.Is(err, account.ErrNotFound) {
			abortWithError(c, http.StatusUnauthorized, account.ErrIncorrectCredentials, "")
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

func (h *authHandler) createToken(c *gin.Context) {
	var r auth.TokenCreateRequest

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

	t, err := h.service.NewToken(
		c.Request.Context(),
		accountId,
		r.Password,
		strings.ToLower(r.Audience),
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - h.service.NewToken: %w", err))

		if errors.Is(err, account.ErrIncorrectPassword) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, auth.TokenCreateResponse{Token: t})
}
