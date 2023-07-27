package media

import (
	"context"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/twirp/common"
)

func (h *Handler) UploadMedia(ctx context.Context, r *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
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
