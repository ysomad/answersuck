package v1

import (
	"context"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ apptwirp.Handler = &MediaHandler{}
	_ pb.MediaService  = &MediaHandler{}
)

type MediaUseCase interface {
	Save(context.Context, entity.Media) (entity.Media, error)
}

type MediaHandler struct {
	media   MediaUseCase
	session *session.Manager
}

func NewMediaHandler(uc MediaUseCase, sm *session.Manager) *MediaHandler {
	return &MediaHandler{
		media:   uc,
		session: sm,
	}
}

func (h *MediaHandler) Handle(m *http.ServeMux) {
	s := pb.NewMediaServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *MediaHandler) UploadMedia(ctx context.Context,
	r *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if r.Url == "" {
		return nil, twirp.RequiredArgumentError("url")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	media, err := entity.NewMedia(r.Url, session.User.ID)
	if err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	answer, err := h.media.Save(ctx, media)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.UploadMediaResponse{
		Media: &pb.Media{
			Url:    answer.URL,
			Type:   pb.MediaType(answer.Type),
			Author: answer.Uploader,
		},
	}, nil
}
