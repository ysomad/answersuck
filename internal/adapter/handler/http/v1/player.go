package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/answersuck/host/internal/domain/player"
	"github.com/answersuck/host/internal/pkg/mime"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type playerHandler struct {
	log      *zap.Logger
	validate validate
	player   playerService
}

func newPlayerMux(d *Deps) *chi.Mux {
	h := playerHandler{
		log:      d.Logger,
		validate: d.Validate,
		player:   d.PlayerService,
	}

	m := chi.NewMux()
	authenticator := mwAuthenticator(d.Logger, &d.Config.Session, d.SessionService)

	m.Get("/{nickname}", h.get)
	m.With(authenticator).Put("/avatar", h.uploadAvatar)

	return m
}

func (h *playerHandler) get(w http.ResponseWriter, r *http.Request) {
	nickname := chi.URLParam(r, "nickname")
	if nickname == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p, err := h.player.GetByNickname(r.Context(), nickname)
	if err != nil {
		h.log.Error("http - v1 - player - get - h.player.GetByNickname", zap.Error(err))

		if errors.Is(err, player.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *playerHandler) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	accountId, err := getAccountId(r.Context())
	if err != nil {
		h.log.Info("http - v1 - player - uploadAvatar - getAccountId", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = r.ParseMultipartForm(player.MaxAvatarSize); err != nil {
		h.log.Info("http - v1 - player - uploadAvatar - r.ParseMultipartForm", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		h.log.Info("http - v1 - player - uploadAvatar - r.FormFile", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	buf := make([]byte, fileHeader.Size)
	if _, err = file.Read(buf); err != nil {
		h.log.Error("http - v1 - player - uploadAvatar - f.Read", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentType, err := mime.NewImageType(http.DetectContentType(buf))
	if err != nil {
		h.log.Info("http - v1 - player - uploadAvatar - mime.NewImageType", zap.Error(err))
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	uuid4, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s-%s", uuid4, fileHeader.Filename)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		h.log.Error("http - v1 - player - uploadAvatar - os.OpenFile", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, bytes.NewReader(buf)); err != nil {
		h.log.Error("http - v1 - player - uploadAvatar - io.Copy", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.player.UploadAvatar(r.Context(), player.UploadAvatarDTO{
		AccountId:   accountId,
		Filename:    filename,
		FileSize:    fileHeader.Size,
		ContentType: contentType,
	}); err != nil {
		h.log.Error("http - v1 - player - setAvatar - h.player.SetAvatar", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
