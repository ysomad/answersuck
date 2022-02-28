package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/domain"
	"github.com/quizlyfun/quizly-backend/internal/service"

	"github.com/quizlyfun/quizly-backend/pkg/logging"
)

func sessionMiddleware(l logging.Logger, cfg *app.Config, session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(cfg.SessionCookie)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - sessionMiddleware - c.Cookie: %w", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s, err := session.GetByID(c.Request.Context(), sid)
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

		c.Set("sid", s.ID)
		c.Set("aid", s.AccountID)
		c.Next()
	}
}

func tokenMiddleware(l logging.Logger, auth service.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		aid, err := accountID(c)
		if err != nil {
			l.Error(fmt.Errorf("http - v1 - middleware - tokenMiddleware - accountID: %w", err))
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
		l.Debug(currAud)

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

		c.Next()
	}
}
