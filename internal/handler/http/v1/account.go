package v1

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"

	"github.com/answersuck/vault/pkg/logging"
)

type AccountService interface {
	Create(ctx context.Context, r account.CreateRequest) (*account.Account, error)
	Delete(ctx context.Context, accountId string) error

	RequestVerification(ctx context.Context, accountId string) error
	Verify(ctx context.Context, code string) error

	ResetPassword(ctx context.Context, login string) error
	SetPassword(ctx context.Context, token, password string) error
}

type accountHandler struct {
	cfg     *config.Aggregate
	log     logging.Logger
	v       ValidationModule
	service AccountService
}

func newAccountRouter(d *Deps) *fiber.App {
	h := &accountHandler{
		cfg:     d.Config,
		log:     d.Logger,
		v:       d.ValidationModule,
		service: d.AccountService,
	}

	r := fiber.New()

	r.Post("/", h.create)

	protected := r.Group("/",
		sessionMW(d.Logger, &d.Config.Session, d.SessionService),
		tokenMW(d.Logger, d.AuthService))
	{
		protected.Delete("/", h.delete)
	}

	verification := r.Group("/verification")
	{
		verification.Post("/",
			sessionMW(d.Logger, &d.Config.Session, d.SessionService),
			h.requestVerification)
		verification.Put("/", h.verify)
	}

	password := r.Group("/password")
	{
		password.Post("/", h.resetPassword)
		password.Put("/", h.setPassword)
	}

	return r
}

func (h *accountHandler) create(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	var r account.CreateRequest

	if err := c.BodyParser(&r); err != nil {
		h.log.Info("http - v1 - account - create - c.BodyParser: %w", err)

		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
	}

	if err := h.v.ValidateStruct(r); err != nil {
		h.log.Info("http - v1 - account - create - ValidateStruct: %w", err)

		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
	}

	_, err := h.service.Create(c.Context(), r)
	if err != nil {
		h.log.Error("http - v1 - account - create - h.service.Create: %w", err)

		switch {

		case errors.Is(err, account.ErrAlreadyExist):
			return errorResp(c, fiber.StatusConflict, account.ErrAlreadyExist, "")

		case errors.Is(err, account.ErrForbiddenNickname):
			return errorResp(c, fiber.StatusBadRequest, account.ErrForbiddenNickname, "")

		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *accountHandler) delete(c *fiber.Ctx) error {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - account - delete - getAccountId: %w", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err = h.service.Delete(c.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - delete - h.service.Delete: %w", err)

		if errors.Is(err, account.ErrNotDeleted) {
			return errorResp(c, fiber.StatusBadRequest, account.ErrAlreadyArchived, "")
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.Session.CookieName,
		Value:    "",
		MaxAge:   -1,
		Secure:   h.cfg.Session.CookieSecure,
		HTTPOnly: h.cfg.Session.CookieHTTPOnly,
	})

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *accountHandler) requestVerification(c *fiber.Ctx) error { return nil }

func (h *accountHandler) verify(c *fiber.Ctx) error        { return nil }
func (h *accountHandler) resetPassword(c *fiber.Ctx) error { return nil }
func (h *accountHandler) setPassword(c *fiber.Ctx) error   { return nil }
