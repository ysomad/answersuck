package v1

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type languageHandler struct {
	log     *zap.Logger
	service LanguageService
}

func newLanguageHandler(d *Deps) *languageHandler {
	return &languageHandler{
		log:     d.Logger,
		service: d.LanguageService,
	}
}

func newLanguageRouter(d *Deps) *fiber.App {
	h := newLanguageHandler(d)
	r := fiber.New()

	r.Get("/", h.getAll)

	return r
}

func (h *languageHandler) getAll(c *fiber.Ctx) error {
	l, err := h.service.GetAll(c.Context())
	if err != nil {
		h.log.Error("http - v1 - language - getAll - h.service.GetAll: %w", zap.Error(err))
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusOK).JSON(listResp{Result: l})
	return nil
}
