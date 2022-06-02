package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/auth"
	"github.com/answersuck/vault/internal/domain/session"
	"github.com/answersuck/vault/pkg/logging"
)

type AuthService interface {
	Login(ctx context.Context, login, password string, d session.Device) (*session.Session, error)

	NewToken(ctx context.Context, accountId, password, audience string) (string, error)
	ParseToken(ctx context.Context, token, audience string) (string, error)
}

type authHandler struct {
	cfg     *config.Aggregate
	log     logging.Logger
	v       ValidationModule
	service AuthService
	session SessionService
}

func newAuthRouter(d *Deps) *fiber.App {
	h := &authHandler{
		cfg:     d.Config,
		log:     d.Logger,
		v:       d.ValidationModule,
		service: d.AuthService,
		session: d.SessionService,
	}

	r := fiber.New()

	r.Post("/login", deviceMW, h.login)

	authenticated := r.Group("/", sessionMW(d.Logger, &d.Config.Session, d.SessionService))
	{
		authenticated.Post("logout", h.logout)
		authenticated.Post("token", h.createToken)
	}

	return r
}

func (h *authHandler) login(c *fiber.Ctx) error {
	if sessionId := c.Cookies(h.cfg.Session.CookieName); sessionId != "" {
		return errorResp(c, fiber.StatusBadRequest, auth.ErrAlreadyLoggedIn, "")
	}

	var r auth.LoginRequest

	if err := c.BodyParser(&r); err != nil {
		h.log.Info("http - v1 - auth - login - c.BodyParser: %w", err)

		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
	}

	if err := h.v.ValidateStruct(r); err != nil {
		h.log.Info("http - v1 - auth - login - h.v.ValidateStruct: %w", err)

		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
	}

	d, err := getDevice(c)
	if err != nil {
		h.log.Error("http - v1 - auth - login - getDevice: %w", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	s, err := h.service.Login(c.Context(), r.Login, r.Password, d)
	if err != nil {
		h.log.Error("http - v1 - auth - login - h.service.Login: %w", err)

		if errors.Is(err, account.ErrIncorrectPassword) || errors.Is(err, account.ErrNotFound) {
			return errorResp(c, fiber.StatusUnauthorized, account.ErrIncorrectCredentials, "")
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.Session.CookieName,
		Value:    s.Id,
		MaxAge:   s.MaxAge,
		Secure:   h.cfg.Session.CookieSecure,
		HTTPOnly: h.cfg.Session.CookieHTTPOnly,
	})

	c.Status(fiber.StatusOK)

	return nil
}

func (h *authHandler) logout(c *fiber.Ctx) error {
	// sessionId, err := getSessionId(c)
	// if err != nil {
	// 	h.log.Info("http - v1 - auth - logout - getSessionId: %w", err)
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// if err := h.session.Terminate(c.Request.Context(), sessionId); err != nil {
	// 	h.log.Error("http - v1 - auth - logout - h.session.Terminate: %w", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// c.SetCookie(h.cfg.Session.CookieName, "", -1, "", "", h.cfg.Session.CookieSecure, h.cfg.Session.CookieHTTPOnly)
	// c.Status(http.StatusNoContent)
	return nil
}

func (h *authHandler) createToken(c *fiber.Ctx) error {
	// var r auth.TokenCreateRequest

	// if err := c.ShouldBindJSON(&r); err != nil {
	// 	h.log.Info(err.Error())
	// 	abortWithError(c, http.StatusBadRequest, errInvalidRequestBody, h.t.TranslateError(err))
	// 	return
	// }

	// accountId, err := getAccountId(c)
	// if err != nil {
	// 	h.log.Error("http - v1 - auth - createToken - getAccountId: %w", err)
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// t, err := h.service.NewToken(
	// 	c.Request.Context(),
	// 	accountId,
	// 	r.Password,
	// 	strings.ToLower(r.Audience),
	// )
	// if err != nil {
	// 	h.log.Error("http - v1 - auth - createToken - h.service.NewToken: %w", err)

	// 	if errors.Is(err, account.ErrIncorrectPassword) {
	// 		c.AbortWithStatus(http.StatusForbidden)
	// 		return
	// 	}

	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// c.JSON(http.StatusOK, auth.TokenCreateResponse{Token: t})
	return nil
}
