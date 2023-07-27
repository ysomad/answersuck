package pack

import (
	"context"
	"errors"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (h *Handler) CreatePack(ctx context.Context, r *pb.CreatePackRequest) (*pb.CreatePackResponse, error) {
	session, ok := appctx.GetSession(ctx)
	if !ok {
		return nil, twirp.Unauthenticated.Error(apperr.MsgUnauthorized)
	}

	if r.PackName == "" {
		return nil, twirp.RequiredArgumentError("pack_name")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	packID, err := h.pack.Save(ctx, &entity.Pack{
		Name:       r.PackName,
		Author:     session.User.ID,
		CoverURL:   r.CoverUrl,
		CreateTime: time.Now(),
	}, r.Tags)
	if err != nil {
		if errors.Is(err, apperr.MediaNotFound) {
			return nil, twirp.InvalidArgument.Error(apperr.MsgPackCoverNotFound)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreatePackResponse{PackId: packID}, nil
}
