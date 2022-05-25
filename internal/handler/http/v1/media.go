package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/answersuck/vault/internal/config"
	"github.com/answersuck/vault/internal/domain/media"
	"github.com/answersuck/vault/pkg/logging"
)

type MediaService interface {
	UploadAndSave(ctx context.Context, dto *media.UploadDTO) (media.Media, error)
}

type mediaHandler struct {
	t       ErrorTranslator
	cfg     *config.Aggregate
	log     logging.Logger
	service MediaService
}

func newMediaHandler(r *gin.RouterGroup, d *Deps) {
	h := &mediaHandler{
		t:       d.ErrorTranslator,
		cfg:     d.Config,
		log:     d.Logger,
		service: d.MediaService,
	}

	media := r.Group("media")
	{
		media.POST("", sessionMiddleware(d.Logger, &d.Config.Session, d.SessionService), h.upload)
	}
}

const (
	mediaFormKey = "media"
)

func (h *mediaHandler) upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, media.MaxUploadSize)

	f, fh, err := c.Request.FormFile(mediaFormKey)
	if err != nil {
		h.log.Error("http - v1 - media - upload - c.Request.FormFile: %w", err)
		abortWithError(c, http.StatusBadRequest, err, "")
		return
	}

	defer f.Close()

	buf := make([]byte, fh.Size)

	if _, err := f.Read(buf); err != nil {
		h.log.Error("http - v1 - media - upload - file.Read: %w", err)
		abortWithError(c, http.StatusBadRequest, err, "")
		return
	}

	accountId, err := getAccountId(c)
	if err != nil {
		h.log.Error("http - v1 - media - upload - getAccountId: %w", err)
		abortWithError(c, http.StatusUnauthorized, err, "")
		return
	}

	m, err := h.service.UploadAndSave(c.Request.Context(), &media.UploadDTO{
		Filename:    fh.Filename,
		ContentType: http.DetectContentType(buf),
		Size:        fh.Size,
		Buf:         buf,
		AccountId:   accountId,
	})

	if err != nil {
		h.log.Error("http - v1 - media - upload: %w", err)

		if errors.Is(err, media.ErrInvalidMimeType) {
			abortWithError(c, http.StatusBadRequest, media.ErrInvalidMimeType, "")
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, m)
}
