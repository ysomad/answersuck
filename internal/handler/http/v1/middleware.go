package v1

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/session"
	"github.com/answersuck/vault/pkg/logging"
)

// sessionMW check if request is authenticated and sets accountId and sessionId to locals (context)
func sessionMW(l logging.Logger, cfg *config.Session,
	s SessionService) fiber.Handler {

	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies(cfg.CookieName)
		if sessionId == "" {
			l.Info("http - v1 - middleware - sessionMW - c.Cookies: cookie '%s' not found", cfg.CookieName)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		sess, err := s.GetById(c.Context(), sessionId)
		if err != nil {
			l.Error("http - v1 - middleware - sessionMW - s.GetById: %w", err)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if sess.IP != c.IP() || sess.UserAgent != c.Get(fiber.HeaderUserAgent) {
			l.Error("http - v1 - middleware - sessionMW: %w", session.ErrDeviceMismatch)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		c.Locals(sessionIdKey, sess.Id)
		c.Locals(accountIdKey, sess.AccountId)

		return c.Next()
	}
}

// verifiedMW is simillar to sessionMW but also checks if account is verified,
// aborts if not.
//
// should be used instead of sessionMW
func verifiedMW(l logging.Logger, cfg *config.Session,
	s SessionService) fiber.Handler {

	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies(cfg.CookieName)
		if sessionId == "" {
			l.Info("http - v1 - middleware - verifiedMW - c.Cookies: cookie '%s' not found", cfg.CookieName)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		d, err := s.GetByIdWithVerified(c.Context(), sessionId)
		if err != nil {
			l.Error("http - v1 - middleware - verifiedMW - s.GetById: %w", err)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if d.Session.IP != c.IP() || d.Session.UserAgent != c.Get(fiber.HeaderUserAgent) {
			l.Error("http - v1 - middleware - verifiedMW: %w", session.ErrDeviceMismatch)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if !d.Verified {
			l.Info("http - v1 - middleware - verifiedMW - !d.Verified: %w", account.ErrNotEnoughRights)

			return errorResp(c, fiber.StatusForbidden, account.ErrNotEnoughRights, "")
		}

		c.Locals(sessionIdKey, d.Session.Id)
		c.Locals(accountIdKey, d.Session.AccountId)

		return c.Next()
	}
}

// tokenMW parses and validates security token
func tokenMW(l logging.Logger, auth AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accountId, err := getAccountId(c)
		if err != nil {
			l.Info("http - v1 - middleware - tokenMW - getAccountId: %w", err)

			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		t := c.Query("token")
		if t == "" {
			l.Info("http - v1 - middleware - tokenMW - c.Query: %w", account.ErrEmptyPasswordResetToken)

			c.Status(fiber.StatusForbidden)
			return nil
		}

		reqURI := strings.ToLower(c.BaseURL() + c.OriginalURL())

		sub, err := auth.ParseToken(c.Context(), t, reqURI)
		if err != nil {
			l.Error("http - v1 - middleware - tokenMW - auth.ParseToken: %w", err)

			c.Status(fiber.StatusForbidden)
			return nil
		}

		if sub != accountId {
			l.Info("http - v1 - middleware - tokenMW: %w", err)

			c.Status(fiber.StatusForbidden)
			return nil
		}

		return c.Next()
	}
}

// deviceMW sets session.Device object to locals
func deviceMW(c *fiber.Ctx) error {
	c.Locals(deviceKey, session.Device{
		IP:        c.IP(),
		UserAgent: c.Get(fiber.HeaderUserAgent),
	})

	return c.Next()
}
