package v1

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/answersuck/host/internal/domain/media"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type mediaHandler struct {
	log      *zap.Logger
	validate validate
	media    mediaService
}

func newMediaHandler(d *Deps) http.Handler {
	h := mediaHandler{
		log:      d.Logger,
		validate: d.Validate,
		media:    d.MediaService,
	}

	r := chi.NewRouter()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	r.With(verificator).Post("/", h.upload)

	return r
}

func (h *mediaHandler) upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	accountId, err := getAccountId(ctx)
	if err != nil {
		h.log.Info("http - v1 - media - upload - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = r.ParseMultipartForm(media.MaxUploadSize); err != nil {
		h.log.Info("http - v1 - media - upload - r.ParseMultipartForm", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := r.FormFile("media")
	if err != nil {
		h.log.Info("http - v1 - media - upload - r.FormFile", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	buf := make([]byte, fileHeader.Size)
	if _, err = file.Read(buf); err != nil {
		h.log.Error("http - v1 - media - upload - f.Read", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := media.New(fileHeader.Filename, accountId, media.Type(http.DetectContentType(buf)))
	if err != nil {
		h.log.Error("http - v1 - media - upload - media.NewMedia", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	f, err := os.OpenFile(m.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		h.log.Error("http - v1 - media - upload - os.OpenFile", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, bytes.NewReader(buf)); err != nil {
		h.log.Error("http - v1 - media - upload - io.Copy", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := h.media.UploadAndSave(ctx, m, fileHeader.Size)
	if err != nil {
		h.log.Error("http - v1 - media - upload - h.media.UploadAndSave", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, resp)
}

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
