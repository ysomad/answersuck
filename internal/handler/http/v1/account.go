package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/account"

	"github.com/answersuck/vault/pkg/logging"
)

type accountHandler struct {
	cfg          *config.Aggregate
	log          logging.Logger
	v            ValidationModule
	account      AccountService
	verification VerificationService
}

func newAccountHandler(d *Deps) *accountHandler {
	return &accountHandler{
		cfg:          d.Config,
		log:          d.Logger,
		v:            d.ValidationModule,
		account:      d.AccountService,
		verification: d.VerificationService,
	}
}

func newAccountRouter(d *Deps) *fiber.App {
	h := newAccountHandler(d)

	r := fiber.New()

	r.Post("/", h.create)

	protected := r.Group("/",
		sessionMW(d.Logger, &d.Config.Session, d.SessionService),
		tokenMW(d.Logger, d.TokenService))
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

	var r account.CreateReq

	if err := c.BodyParser(&r); err != nil {
		h.log.Info("http - v1 - account - create - c.BodyParser: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
	}

	if err := h.v.ValidateStruct(r); err != nil {
		h.log.Info("http - v1 - account - create - ValidateStruct: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
	}

	_, err := h.account.Create(c.Context(), r)
	if err != nil {
		h.log.Error("http - v1 - account - create - h.account.Create: %w", err)

		switch {
		case errors.Is(err, account.ErrAlreadyExist):
			return errorResp(c, fiber.StatusConflict, account.ErrAlreadyExist, "")
		case errors.Is(err, account.ErrForbiddenNickname):
			return errorResp(c, fiber.StatusBadRequest, account.ErrForbiddenNickname, "")
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *accountHandler) delete(c *fiber.Ctx) error {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - account - delete - getAccountId: %w", err)

		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	if err = h.account.Delete(c.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - delete - h.account.Delete: %w", err)

		if errors.Is(err, account.ErrNotDeleted) {
			return errorResp(c, fiber.StatusBadRequest, account.ErrAlreadyArchived, "")
		}

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

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *accountHandler) requestVerification(c *fiber.Ctx) error {
	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - account - requestVerification - getAccountId: %w", err)

		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	if err = h.verification.Request(c.Context(), accountId); err != nil {
		h.log.Error("http - v1 - account - requestVerification - h.verification.Request")

		if errors.Is(err, account.ErrAlreadyVerified) {
			return errorResp(c, fiber.StatusBadRequest, account.ErrAlreadyVerified, "")
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusAccepted)
	return nil
}

func (h *accountHandler) verify(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return errorResp(c, fiber.StatusBadRequest, account.ErrEmptyVerificationCode, "")
	}

	if err := h.verification.Verify(c.Context(), code); err != nil {
		h.log.Error("http - v1 - account - verify - h.verification.Verify: %w", err)

		if errors.Is(err, account.ErrAlreadyVerified) {
			return errorResp(c, fiber.StatusNotFound, account.ErrAlreadyVerified, "")
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func (h *accountHandler) resetPassword(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	var r account.ResetPasswordReq

	if err := c.BodyParser(&r); err != nil {
		h.log.Info("http - v1 - account - resetPassword - c.BodyParser: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
	}

	if err := h.v.ValidateStruct(r); err != nil {
		h.log.Info("http - v1 - account - resetPassword - ValidateStruct: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
	}

	if err := h.account.ResetPassword(c.Context(), r.Login); err != nil {
		h.log.Error("http - v1 - account - resetPassword - h.password.Reset: %w", err)

		if errors.Is(err, account.ErrNotFound) {
			return errorResp(c, fiber.StatusNotFound, account.ErrNotFound, "")
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusAccepted)
	return nil
}

func (h *accountHandler) setPassword(c *fiber.Ctx) error {
	t := c.Query("token")
	if t == "" {
		return errorResp(c, fiber.StatusBadRequest, account.ErrEmptyPasswordResetToken, "")
	}

	c.Accepts(fiber.MIMEApplicationJSON)

	var r account.SetPasswordReq

	if err := c.BodyParser(&r); err != nil {
		h.log.Info("http - v1 - account - setPassword - c.BodyParser: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
	}

	if err := h.v.ValidateStruct(r); err != nil {
		h.log.Info("http - v1 - account - setPassword - ValidateStruct: %w", err)
		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
	}

	if err := h.account.SetPassword(c.Context(), t, r.Password); err != nil {
		h.log.Error("http - v1 - account - setPassword - h.account.SetPassword: %w", err)

		if errors.Is(err, account.ErrPasswordResetTokenExpired) ||
			errors.Is(err, account.ErrPasswordTokenNotFound) {

			c.Status(fiber.StatusForbidden)
			return nil
		}

		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
