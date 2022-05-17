package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

const userAgentHeader = "User-Agent"

// sessionMiddleware looking for a cookie with session id, sets account id and session id to context
func sessionMiddleware(l logging.Logger, cfg *config.Session, service SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie(cfg.CookieKey)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s, err := service.GetById(c.Request.Context(), sessionId)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - s.Get: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if s.IP != c.ClientIP() || s.UserAgent != c.GetHeader(userAgentHeader) {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware: %w", session.ErrDeviceMismatch))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		l.Info(s.AccountId)

		c.Set(sessionIdKey, s.Id)
		c.Set(accountIdKey, s.AccountId)
		c.Next()
	}
}

// tokenMiddleware parses and validates security token
func tokenMiddleware(l logging.Logger, auth AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountId, err := getAccountId(c)
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

		if sub != accountId {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware: %w", err))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set(audienceKey, currAud)
		c.Next()
	}
}

func deviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := session.Device{
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader(userAgentHeader),
		}

		c.Set(deviceKey, d)
		c.Next()
	}
}
