package v1

import "go.uber.org/zap"

type topicHandler struct {
	log     *zap.Logger
	v       ValidationModule
	service TopicService
}

// func newTopicHandler(d *Deps) *topicHandler {
// 	return &topicHandler{
// 		log:     d.Logger,
// 		v:       d.ValidationModule,
// 		service: d.TopicService,
// 	}
// }
//
// func newTopicRouter(d *Deps) *fiber.App {
// 	h := newTopicHandler(d)
// 	r := fiber.New()
//
// 	r.Get("/", h.getAll)
// 	r.Post("/",
// 		mwVerificator(d.Logger, &d.Config.Session, d.SessionService),
// 		h.create)
//
// 	return r
// }
//
// func (h *topicHandler) create(c *fiber.Ctx) error {
// 	var r topic.CreateReq
//
// 	if err := c.BodyParser(&r); err != nil {
// 		h.log.Info("http - v1 - topic - create - c.BodyParser: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
// 	}
//
// 	if err := h.v.ValidateStruct(r); err != nil {
// 		h.log.Info("http - v1 - topic - create - ValidateStruct: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
// 	}
//
// 	t, err := h.service.Create(c.Context(), r)
// 	if err != nil {
// 		h.log.Error("http - v1 - topic - create - h.service.Create: %w", err)
//
// 		if errors.Is(err, topic.ErrLanguageNotFound) {
// 			return errorResp(c, fiber.StatusBadRequest, topic.ErrLanguageNotFound, "")
// 		}
//
// 		c.Status(fiber.StatusInternalServerError)
// 		return nil
// 	}
//
// 	c.Status(fiber.StatusOK).JSON(t)
// 	return nil
// }
//
// func (h *topicHandler) getAll(c *fiber.Ctx) error {
// 	t, err := h.service.GetAll(c.Context())
// 	if err != nil {
// 		h.log.Error("http - v1 - topic - getAll - h.service.GetAll: %w", err)
// 		c.Status(fiber.StatusInternalServerError)
// 		return nil
// 	}
//
// 	c.Status(fiber.StatusOK).JSON(listResp{Result: t})
// 	return nil
// }
