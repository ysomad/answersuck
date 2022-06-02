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
func sessionMiddleware(l logging.Logger, cfg *config.Session, s SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie(cfg.CookieName)
		if err != nil {
			l.Error("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess, err := s.GetById(c.Request.Context(), sessionId)
		if err != nil {
			l.Error("http - v1 - middleware - sessionMiddleware - s.GetById: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if sess.IP != c.ClientIP() || sess.UserAgent != c.GetHeader(userAgentHeader) {
			l.Error("http - v1 - middleware - sessionMiddleware: %w", session.ErrDeviceMismatch)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(sessionIdKey, sess.Id)
		c.Set(accountIdKey, sess.AccountId)
		c.Next()
	}
}

// protectionMiddleware checks if account is verified before
func protectionMiddleware(l logging.Logger, cfg *config.Session, s SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie(cfg.CookieName)
		if err != nil {
			l.Error("http - v1 - middleware - protectionMiddleware - c.Cookie: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess, err := s.GetByIdWithVerified(c.Request.Context(), sessionId)
		if err != nil {
			l.Error("http - v1 - middleware - protectionMiddleware - s.Get: %w", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if sess.Session.IP != c.ClientIP() || sess.Session.UserAgent != c.GetHeader(userAgentHeader) {
			l.Error("http - v1 - middleware - protectionMiddleware: %w", session.ErrDeviceMismatch)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !sess.Verified {
			l.Error("http - v1 - middleware - protectionMiddleware - getAccountVerified: %w", err)
			abortWithError(c, http.StatusForbidden, account.ErrNotEnoughRights, "")
			return
		}

		c.Set(sessionIdKey, sess.Session.Id)
		c.Set(accountIdKey, sess.Session.AccountId)
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
