package v1

import (
	"go.uber.org/zap"

	"github.com/answersuck/vault/internal/config"
)

type mediaHandler struct {
	cfg     *config.Aggregate
	log     *zap.Logger
	v       ValidationModule
	service MediaService
}

// func newMediaHandler(d *Deps) *mediaHandler {
// 	h := mediaHandler{
// 		cfg:     d.Config,
// 		log:     d.Logger,
// 		v:       d.ValidationModule,
// 		service: d.MediaService,
// 	}
//
// 	return &h
// }
//
// func newMediaRouter(d *Deps) *fiber.App {
// 	h := newMediaHandler(d)
// 	r := fiber.New()
//
// 	r.Post("/", mwVerificator(d.Logger, &d.Config.Session, d.SessionService), h.upload)
//
// 	return r
// }
//
// func (h *mediaHandler) upload(c *fiber.Ctx) error {
// 	c.Accepts(fiber.MIMEMultipartForm)
//
// 	accountId, err := getAccountId(c)
// 	if err != nil {
// 		h.log.Info("http - v1 - media - upload - getAccountId: %w", err)
// 		c.Status(fiber.StatusUnauthorized)
// 		return nil
// 	}
//
// 	fh, err := c.FormFile("media")
// 	if err != nil {
// 		h.log.Info("http - v1 - media - upload - c.FormFile: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, err, "")
// 	}
//
// 	f, err := fh.Open()
// 	if err != nil {
// 		h.log.Info("http - v1 - media - upload - c.FormFile: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, err, "")
// 	}
//
// 	defer f.Close()
//
// 	buf := make([]byte, fh.Size)
//
// 	if _, err := f.Read(buf); err != nil {
// 		h.log.Info("http - v1 - media - upload - c.FormFile: %w", err)
// 		return errorResp(c, fiber.StatusBadRequest, err, "")
// 	}
//
// 	//c.SaveFile(f, fmt.Sprintf("./static/media/%s", f.Filename))
//
// 	m, err := h.service.UploadAndSave(c.Context(), &media.UploadDTO{
// 		Filename:    fh.Filename,
// 		ContentType: http.DetectContentType(buf),
// 		Size:        fh.Size,
// 		Buf:         buf,
// 		AccountId:   accountId,
// 	})
// 	if err != nil {
// 		h.log.Error("http - v1 - media - upload: %w", err)
//
// 		if errors.Is(err, media.ErrInvalidMimeType) {
// 			return errorResp(c, fiber.StatusBadRequest, media.ErrInvalidMimeType, "")
// 		}
//
// 		c.Status(fiber.StatusInternalServerError)
// 		return nil
// 	}
//
// 	c.Status(fiber.StatusOK).JSON(m)
// 	return nil
// }
