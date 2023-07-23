package v1

import (
	"context"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/media/v1"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) SaveMedia(ctx context.Context, p *pb.SaveMediaRequest) (*pb.SaveMediaResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if p.Url == "" {
		return nil, twirp.RequiredArgumentError("url")
	}

	media, err := entity.NewMedia(p.Url, session.Player.Nickname)
	if err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	answer, err := h.media.Save(ctx, media)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.SaveMediaResponse{
		Media: &pb.Media{
			Url:        answer.URL,
			Type:       pb.MediaType(answer.Type),
			Author:     answer.Author,
			CreateTime: timestamppb.New(answer.CreateTime),
		},
	}, nil
}
