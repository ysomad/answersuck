package v1

import "go.uber.org/zap"

type tagHandler struct {
	log     *zap.Logger
	v       ValidationModule
	service TagService
}

// func newTagHandler(d *Deps) *tagHandler {
// 	return &tagHandler{
// 		log:     d.Logger,
// 		v:       d.ValidationModule,
// 		service: d.TagService,
// 	}
// }
//
// func newTagRouter(d *Deps) *fiber.App {
// 	h := newTagHandler(d)
// 	r := fiber.New()
//
// 	r.Get("/", h.getAll)
// 	r.Post("/",
// 		mwVerificator(d.Logger, &d.Config.Session, d.SessionService),
// 		h.createMultiple)
//
// 	return r
// }
//
// func (h *tagHandler) createMultiple(c *fiber.Ctx) error {
// 	var r tag.CreateMultipleReq
//
// 	if err := c.BodyParser(&r); err != nil {
// 		h.log.Info("http - v1 - tag - createMultiple - c.BodyParser: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
// 	}
//
// 	if err := h.v.ValidateStruct(r); err != nil {
// 		h.log.Info("http - v1 - tag - createMultiple - ValidateStruct: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
// 	}
//
// 	t, err := h.service.CreateMultiple(c.Context(), r.Tags)
// 	if err != nil {
// 		h.log.Error("http - v1 - tag - createMultiple - h.service.CreateMultiple: %w", err)
// 		c.Status(fiber.StatusInternalServerError)
// 		return nil
// 	}
//
// 	c.Status(fiber.StatusOK).JSON(listResp{Result: t})
// 	return nil
// }
//
// func (h *tagHandler) getAll(c *fiber.Ctx) error {
// 	t, err := h.service.GetAll(c.Context())
// 	if err != nil {
// 		h.log.Error("http - v1 - tag - getAll - h.service.GetAll: %w", err)
// 		c.Status(fiber.StatusInternalServerError)
// 		return nil
// 	}
//
// 	c.Status(fiber.StatusOK).JSON(listResp{Result: t})
// 	return nil
// }
