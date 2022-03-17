package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/internal/service"

	"github.com/quizlyfun/quizly-backend/pkg/logging"
	"github.com/quizlyfun/quizly-backend/pkg/validation"
)

type authHandler struct {
	validation.ErrorTranslator
	cfg  *app.Config
	log  logging.Logger
	auth service.Auth
}

func newAuthHandler(handler *gin.RouterGroup, d *Deps) {
	h := &authHandler{
		d.ErrorTranslator,
		d.Config,
		d.Logger,
		d.AuthService,
	}

	g := handler.Group("auth")
	{
		authenticated := g.Group("", sessionMiddleware(d.Logger, d.Config, d.SessionService))
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
		abortWithValidationError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.TranslateError(err))
	}

	s, err := h.auth.Login(
		c.Request.Context(),
		r.Login,
		r.Password,
		domain.Device{
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		},
	)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - login: %w", err))

		if errors.Is(err, domain.ErrAccountIncorrectPassword) || errors.Is(err, domain.ErrAccountNotFound) {
			abortWithError(c, http.StatusUnauthorized, domain.ErrAccountIncorrectCredentials, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.SessionCookie, s.ID, s.TTL, "", "", h.cfg.CookieSecure, h.cfg.CookieHTTPOnly)
	c.Status(http.StatusOK)
}

func (h *authHandler) logout(c *gin.Context) {
	sid, err := sessionID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - logout - sessionID: %w", err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := h.auth.Logout(c.Request.Context(), sid); err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - logout: %w", err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.SessionCookie, "", -1, "", "", h.cfg.CookieSecure, h.cfg.CookieHTTPOnly)
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
		abortWithValidationError(c, http.StatusBadRequest, ErrInvalidRequestBody, h.TranslateError(err))
		return
	}

	aid, err := accountID(c)
	if err != nil {
		h.log.Error(fmt.Errorf("http - v1 - auth - token - accountID: %w", err))
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
