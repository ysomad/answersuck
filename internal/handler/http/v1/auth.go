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
	repository "github.com/answersuck/vault/internal/repository/psql"

	"github.com/answersuck/vault/pkg/logging"
)

type authService interface {
	Login(ctx context.Context, login, password string, d domain.Device) (*domain.Session, error)
	Logout(ctx context.Context, sid string) error

	NewToken(ctx context.Context, aid, password, audience string) (string, error)
	ParseToken(ctx context.Context, token, audience string) (string, error)
}

type authHandler struct {
	t       errorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	service authService
}

func newAuthHandler(handler *gin.RouterGroup, d *Deps) {
	h := &authHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.AuthService,
	}

	g := handler.Group("auth")
	{
		authenticated := g.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
		{
			authenticated.POST("logout", h.logout)
			authenticated.POST("token", h.tokenCreate)
		}

		g.POST("login", deviceMiddleware(), h.login)
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

	s, err := h.service.Login(
		c.Request.Context(),
		r.Login,
		r.Password,
		d,
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - login - h.service.Login: %w", err))

		if errors.Is(err, domain.ErrAccountIncorrectPassword) || errors.Is(err, repository.ErrNotFound) {
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
	sid := getSessionId(c)

	if err := h.service.Logout(c.Request.Context(), sid); err != nil {
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

	aid, err := getAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - GetAccountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	t, err := h.service.NewToken(
		c.Request.Context(),
		aid,
		r.Password,
		strings.ToLower(r.Audience),
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - h.service.NewToken: %w", err))

		if errors.Is(err, repository.ErrNotFound) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if errors.Is(err, domain.ErrAccountIncorrectPassword) {
			abortWithError(c, http.StatusUnauthorized, domain.ErrAccountIncorrectPassword, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dto.TokenCreateResponse{Token: t})
}
