package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/answersuck/vault/internal/service/repository"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/dto"
	"github.com/answersuck/vault/internal/service"

	"github.com/answersuck/vault/pkg/logging"
	"github.com/answersuck/vault/pkg/validation"
)

type authHandler struct {
	t    validation.ErrorTranslator
	cfg  *config.Aggregate
	log  logging.Logger
	auth service.Auth
}

func newAuthHandler(handler *gin.RouterGroup, d *Deps) {
	h := &authHandler{
		t:    d.ErrorTranslator,
		cfg:  d.Config,
		log:  d.Logger,
		auth: d.AuthService,
	}

	g := handler.Group("auth")
	{
		authenticated := g.Group("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService))
		{
			authenticated.POST("logout", h.logout)
			authenticated.POST("token", h.token)
		}

		g.POST("login", h.login)
	}
}

type loginRequest struct {
	Login    string `json:"login" binding:"required,email|alphanum"`
	Password string `json:"password" binding:"required"`
}

func (h *authHandler) login(c *gin.Context) {
	var r loginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	s, err := h.auth.Login(
		c.Request.Context(),
		r.Login,
		r.Password,
		dto.Device{
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - login: %w", err))

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
	sid := GetSessionId(c)

	if err := h.auth.Logout(c.Request.Context(), sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - logout: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.Session.CookieKey, "", -1, "", "", h.cfg.Cookie.Secure, h.cfg.Cookie.HTTPOnly)
	c.Status(http.StatusNoContent)
}

type tokenRequest struct {
	Audience string `json:"audience" binding:"required,uri"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (h *authHandler) token(c *gin.Context) {
	var r tokenRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		h.log.Info(err.Error())
		abortWithError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.t.TranslateError(err))
		return
	}

	aid, err := GetAccountId(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - accountId: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	t, err := h.auth.NewAccessToken(
		c.Request.Context(),
		aid,
		r.Password,
		strings.ToLower(r.Audience),
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token: %w", err))

		if errors.Is(err, domain.ErrAccountIncorrectPassword) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{t})
}
