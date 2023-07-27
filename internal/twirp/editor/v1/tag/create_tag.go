package tag

import (
	"context"
	"errors"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateTag(ctx context.Context, r *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {
	session, ok := appctx.GetSession(ctx)
	if !ok {
		return nil, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
	}

	if !session.User.Verified {
		return nil, twirp.PermissionDenied.Error(apperr.MsgPlayerNotVerified)
	}

	if r.TagName == "" {
		return nil, twirp.RequiredArgumentError("tag_name")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	tag := entity.Tag{
		Name:       r.TagName,
		Author:     session.User.UserID,
		CreateTime: time.Now(),
	}

	if err := h.tag.Save(ctx, tag); err != nil {
		if errors.Is(err, apperr.ErrTagAlreadyExists) {
			return nil, twirp.AlreadyExists.Error(apperr.MsgTagAlreadyExists)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateTagResponse{
		Tag: &pb.Tag{
			Name:       tag.Name,
			Author:     tag.Author,
			CreateTime: timestamppb.New(tag.CreateTime),
		},
	}, nil
}
