package v1

import "go.uber.org/zap"

type answerHandler struct {
	log     *zap.Logger
	v       ValidationModule
	service AnswerService
}

// func newAnswerHandler(d *Deps) *answerHandler {
// 	return &answerHandler{
// 		log:     d.Logger,
// 		v:       d.ValidationModule,
// 		service: d.AnswerService,
// 	}
// }
//
// func newAnswerRouter(d *Deps) *fiber.App {
// 	h := newAnswerHandler(d)
// 	r := fiber.New()
//
// 	r.Post("/",
// 		mwVerificator(d.Logger, &d.Config.Session, d.SessionService),
// 		h.create)
//
// 	return r
// }
//
// func (h *answerHandler) create(c *fiber.Ctx) error {
// 	var r answer.CreateReq
//
// 	if err := c.BodyParser(&r); err != nil {
// 		h.log.Info("http - v1 - answer - create - c.BodyParser: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
// 	}
//
// 	if err := h.v.ValidateStruct(r); err != nil {
// 		h.log.Info("http - v1 - answer - create - ValidateStruct: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
// 	}
//
// 	a, err := h.service.Create(c.Context(), r)
// 	if err != nil {
// 		h.log.Error("http - v1 - answer - create - h.service.Create: %w", err)
//
// 		switch {
// 		case errors.Is(err, answer.ErrMimeTypeNotAllowed):
// 			return errorResp(c, fiber.StatusBadRequest, answer.ErrMimeTypeNotAllowed, "")
// 		case errors.Is(err, media.ErrNotFound):
// 			return errorResp(c, fiber.StatusBadRequest, answer.ErrMediaNotFound, "")
// 		}
//
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}
//
// 	return c.Status(fiber.StatusOK).JSON(a)
// }
