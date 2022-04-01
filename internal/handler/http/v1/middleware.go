package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/internal/service"

	"github.com/answersuck/vault/pkg/logging"
)

// sessionMiddleware looking for a cookie with session id, sets account id and session id to context
func sessionMiddleware(l logging.Logger, cfg *config.Session, session service.Session) gin.HandlerFunc {
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

		if s.Device.IP != c.ClientIP() || s.Device.UserAgent != c.GetHeader("User-Agent") {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware: %w", domain.ErrSessionDeviceMismatch))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("sid", s.Id)
		c.Set("aid", s.AccountId)
		c.Next()
	}
}

// tokenMiddleware parses and validates access token
func tokenMiddleware(l logging.Logger, auth service.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := accountId(c)
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

		sub, err := auth.ParseAccessToken(c.Request.Context(), t, currAud)
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
