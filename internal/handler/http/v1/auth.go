package v1

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"

	"github.com/answersuck/vault/internal/domain/account"
	"github.com/answersuck/vault/internal/domain/auth"

	"github.com/answersuck/vault/pkg/logging"
)

type authHandler struct {
	cfg     *config.Aggregate
	log     logging.Logger
	v       ValidationModule
	service LoginService
	token   TokenService
	session SessionService
}

func newAuthHandler(d *Deps) *authHandler {
	return &authHandler{
		cfg:     d.Config,
		log:     d.Logger,
		v:       d.ValidationModule,
		service: d.LoginService,
		token:   d.TokenService,
		session: d.SessionService,
	}
}

func newAuthRouter(d *Deps) *fiber.App {
	h := newAuthHandler(d)

	r := fiber.New()

	r.Post("/login", deviceMW, h.login)

	authenticated := r.Group("/", sessionMW(d.Logger, &d.Config.Session, d.SessionService))
	authenticated.Post("logout", h.logout)
	authenticated.Post("token", h.createToken)

	return r
}

func (h *authHandler) login(c *fiber.Ctx) error {
	if sessionId := c.Cookies(h.cfg.Session.CookieName); sessionId != "" {
		return errorResp(c, fiber.StatusBadRequest, auth.ErrAlreadyLoggedIn, "")
	}

	var r auth.LoginReq

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
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	s, err := h.service.Login(c.Context(), r.Login, r.Password, d)
	if err != nil {
		h.log.Error("http - v1 - auth - login - h.service.Login: %w", err)

		if errors.Is(err, account.ErrIncorrectPassword) ||
			errors.Is(err, account.ErrNotFound) {

			return errorResp(c, fiber.StatusUnauthorized, account.ErrIncorrectCredentials, "")
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
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
	sessionId, err := getSessionId(c)
	if err != nil {
		h.log.Info("http - v1 - auth - logout - getSessionId: %w", err)
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	if err := h.session.Terminate(c.Context(), sessionId); err != nil {
		h.log.Error("http - v1 - auth - logout - h.session.Terminate: %w", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.Session.CookieName,
		Value:    "",
		MaxAge:   -1,
		Secure:   h.cfg.Session.CookieSecure,
		HTTPOnly: h.cfg.Session.CookieHTTPOnly,
	})

	c.Status(fiber.StatusOK)
	return nil
}

func (h *authHandler) createToken(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	var r auth.TokenCreateReq

	if err := c.BodyParser(&r); err != nil {
		h.log.Info("http - v1 - auth - createToken - c.BodyParser: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
	}

	if err := h.v.ValidateStruct(r); err != nil {
		h.log.Info("http - v1 - auth - createToken - h.v.ValidateStruct: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
	}

	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - auth - createToken - getAccountId: %w", err)
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	t, err := h.token.Create(c.Context(), auth.TokenCreateDTO{
		AccountId: accountId,
		Password:  r.Password,
		Audience:  strings.ToLower(r.Audience),
	})
	if err != nil {
		h.log.Error("http - v1 - auth - createToken - h.token.Create: %w", err)

		if errors.Is(err, account.ErrIncorrectPassword) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusOK).JSON(auth.TokenCreateResp{Token: t})
	return nil
}
