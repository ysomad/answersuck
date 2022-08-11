package v1

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/answersuck/host/internal/domain/media"
	"github.com/answersuck/host/internal/pkg/mime"
)

type mediaHandler struct {
	log      *zap.Logger
	validate validate
	media    mediaService
}

func newMediaMux(d *Deps) *chi.Mux {
	h := mediaHandler{
		log:      d.Logger,
		validate: d.Validate,
		media:    d.MediaService,
	}

	m := chi.NewMux()
	verificator := mwVerificator(d.Logger, &d.Config.Session, d.SessionService)

	m.With(verificator).Post("/", h.upload)

	return m
}

func (h *mediaHandler) upload(w http.ResponseWriter, r *http.Request) {
	accountId, err := getAccountId(r.Context())
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

	mediaType, err := mime.NewType(http.DetectContentType(buf))
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	m, err := media.New(fileHeader.Filename, accountId, mediaType)
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

	resp, err := h.media.UploadAndSave(r.Context(), m, fileHeader.Size)
	if err != nil {
		h.log.Error("http - v1 - media - upload - h.media.UploadAndSave", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
