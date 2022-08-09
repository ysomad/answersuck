package v1

import (
	"errors"
	"net/http"

	"github.com/answersuck/host/internal/domain/question"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type questionHandler struct {
	log      *zap.Logger
	validate validate
	question questionService
}

func newQuestionMux(d *Deps) *chi.Mux {
	h := questionHandler{
		log:      d.Logger,
		validate: d.Validate,
		question: d.QuestionService,
	}

	m := chi.NewMux()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	m.With(verificator).Post("/", h.create)

	return m
}

type questionCreateReq struct {
	Text       string `json:"text" validate:"required,gte=1,lte=200"`
	AnswerId   uint32 `json:"answer_id" validate:"required"`
	MediaId    string `json:"media_id" validate:"omitempty,uuid4"`
	LanguageId uint8  `json:"language_id" validate:"required"`
}

type questionCreateRes struct {
	Id uint32 `json:"id"`
}

func (h *questionHandler) create(w http.ResponseWriter, r *http.Request) {
	var req questionCreateReq
	if err := h.validate.RequestBody(r.Body, &req); err != nil {
		h.log.Info("http - v1 - question - create - h.validate.RequestBody", zap.Error(err))
		writeValidationErr(w, http.StatusBadRequest, errInvalidRequestBody, h.validate.TranslateError(err))
		return
	}

	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Error("http - v1 - question - create - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	questionId, err := h.question.Create(r.Context(), question.CreateDTO{
		Text:       req.Text,
		AnswerId:   req.AnswerId,
		AccountId:  accountId,
		LanguageId: req.LanguageId,
	})
	if err != nil {
		h.log.Error("http - v1 - question - create - h.question.Create", zap.Error(err))

		if errors.Is(err, question.ErrForeignKeyViolation) {
			writeErr(w, http.StatusBadRequest, question.ErrForeignKeyViolation)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, questionCreateRes{questionId})
}

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
