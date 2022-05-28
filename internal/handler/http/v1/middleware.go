package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"

	"github.com/answersuck/vault/pkg/logging"
)

const userAgentHeader = "User-Agent"

// sessionMiddleware looking for a cookie with session id, sets account id and session id to context
func sessionMiddleware(l logging.Logger, cfg *config.Session, service SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie(cfg.CookieKey)
		if err != nil {
			l.Error("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s, err := service.GetByIdWithVerified(c.Request.Context(), sessionId)
		if err != nil {
			l.Error("http - v1 - middleware - sessionMiddleware - s.Get: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if s.Session.IP != c.ClientIP() || s.Session.UserAgent != c.GetHeader(userAgentHeader) {
			l.Error("http - v1 - middleware - sessionMiddleware: %w", session.ErrDeviceMismatch)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(sessionIdKey, s.Session.Id)
		c.Set(accountIdKey, s.Session.AccountId)
		c.Set(accountVerifiedKey, s.AccountVerified)
		c.Next()
	}
}

// protectionMiddleware checks if account is verified from context
func protectionMiddleware(l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		a, err := getAccountVerified(c)
		if err != nil || !a {
			l.Error("http - v1 - middleware - protectionMiddleware - getAccountVerified: %w", err)
			abortWithError(c, http.StatusForbidden, account.ErrNotEnoughRights, "")
			return
		}

		c.Next()
	}
}

// tokenMiddleware parses and validates security token
func tokenMiddleware(l logging.Logger, auth AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountId, err := getAccountId(c)
		if err != nil {
			l.Error("http - v1 - middleware - tokenMiddleware - accountId: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t, found := c.GetQuery("token")
		if !found || t == "" {
			l.Error("http - v1 - middleware - tokenMiddleware - c.GetQuery: %w", err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		currAud := strings.ToLower(c.Request.Host + c.FullPath())

		sub, err := auth.ParseToken(c.Request.Context(), t, currAud)
		if err != nil {
			l.Error("http - v1 - middleware - tokenMiddleware - auth.ParseAccessToken: %w", err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if sub != accountId {
			l.Error("http - v1 - middleware - tokenMiddleware: %w", err)
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
