package v1

import "go.uber.org/zap"

type questionHandler struct {
	log     *zap.Logger
	v       ValidationModule
	service QuestionService
}

// func newQuestionHandler(d *Deps) *questionHandler {
// 	return &questionHandler{
// 		log:     d.Logger,
// 		v:       d.ValidationModule,
// 		service: d.QuestionService,
// 	}
// }
//
// func newQuestionRouter(d *Deps) *fiber.App {
// 	h := newQuestionHandler(d)
// 	r := fiber.New()
//
// 	r.Get("/", h.getAll)
// 	r.Post("/", mwVerificator(d.Logger, &d.Config.Session, d.SessionService), h.create)
// 	r.Get(":questionId", h.getById)
//
// 	return r
// }
//
// func (h *questionHandler) create(c *fiber.Ctx) error {
// 	var r question.CreateReq
//
// 	if err := c.BodyParser(&r); err != nil {
// 		h.log.Info("http - v1 - question - create - c.BodyParser: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, err.Error())
// 	}
//
// 	if err := h.v.ValidateStruct(r); err != nil {
// 		h.log.Info("http - v1 - question - create - ValidateStruct: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, errInvalidRequestBody, h.v.TranslateError(err))
// 	}
//
// 	accountId, err := getAccountId(c)
// 	if err != nil {
// 		h.log.Error("http - v1 - question - create - getAccountId: %w", err)
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	}
//
// 	q := question.Question{
// 		Text:       r.Text,
// 		AnswerId:   r.AnswerId,
// 		AccountId:  accountId,
// 		LanguageId: r.LanguageId,
// 	}
//
// 	if r.MediaId != "" {
// 		q.MediaId = &r.MediaId
// 	}
//
// 	res, err := h.service.Create(c.Context(), &q)
// 	if err != nil {
// 		h.log.Error("http - v1 - question - create - h.service.Create :%w", err)
//
// 		if errors.Is(err, question.ErrForeignKeyViolation) {
// 			return errorResp(c, fiber.StatusBadRequest, question.ErrForeignKeyViolation, "")
// 		}
//
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}
//
// 	return c.Status(fiber.StatusOK).JSON(res)
// }
//
// func (h *questionHandler) getAll(c *fiber.Ctx) error {
// 	qs, err := h.service.GetAll(c.Context())
// 	if err != nil {
// 		h.log.Error("http - v1 - question - getAll - h.service.GetAll: %w", err)
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}
//
// 	return c.Status(fiber.StatusOK).JSON(listResp{Result: qs})
// }
//
// func (h *questionHandler) getById(c *fiber.Ctx) error {
// 	s := c.Params("questionId")
//
// 	questionId, err := strconv.Atoi(s)
// 	if err != nil {
// 		h.log.Error("http - v1 - question - getById: %w", err)
// 		return c.SendStatus(fiber.StatusNotFound)
// 	}
//
// 	q, err := h.service.GetById(c.Context(), questionId)
// 	if err != nil {
// 		h.log.Error("http - v1 - question - getById: %w", err)
//
// 		if errors.Is(err, question.ErrNotFound) {
// 			return errorResp(c, fiber.StatusNotFound, question.ErrNotFound, "")
// 		}
//
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}
//
// 	return c.Status(fiber.StatusOK).JSON(q)
// }
