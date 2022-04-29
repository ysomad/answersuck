package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"

	"github.com/answersuck/vault/pkg/logging"
)

// sessionMiddleware looking for a cookie with session id, sets account id and session id to context
func sessionMiddleware(l logging.Logger, cfg *config.Session, session sessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(cfg.CookieKey)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s, err := session.GetById(c.Request.Context(), sid)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - s.Get: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if s.IP != c.ClientIP() || s.UserAgent != c.GetHeader("User-Agent") {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware: %w", domain.ErrSessionDeviceMismatch))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		l.Info(s.AccountId)

		c.Set("sid", s.Id)
		c.Set("aid", s.AccountId)
		c.Next()
	}
}

const accountParam = "accountId"

// accountParamMiddleware checks if account id from context and query param are the same
func accountParamMiddleware(l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := getAccountId(c)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - accountParamMiddleware - getAccountId: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if aid != c.Param(accountParam) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// tokenMiddleware parses and validates access token
func tokenMiddleware(l logging.Logger, auth authService) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := getAccountId(c)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - accountId: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t, found := c.GetQuery("token")
		if !found || t == "" {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - c.GetQuery: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		currAud := strings.ToLower(c.Request.Host + c.FullPath())

		sub, err := auth.ParseToken(c.Request.Context(), t, currAud)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - auth.ParseAccessToken: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if sub != aid {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("aud", currAud)
		c.Next()
	}
}

const uaHeader = "User-Agent"

func deviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := domain.Device{
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader(uaHeader),
		}

		c.Set(deviceKey, d)
		c.Next()
	}
}
