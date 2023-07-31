package round

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) AddTopic(ctx context.Context, r *pb.AddTopicRequest) (*emptypb.Empty, error) {
	if _, err := common.CheckPlayerVerification(ctx); err != nil {
		return nil, err
	}

	if r.RoundId == 0 {
		return nil, twirp.RequiredArgumentError("round_id")
	}

	if r.TopicId == 0 {
		return nil, twirp.RequiredArgumentError("topic_id")
	}

	if err := h.round.AddTopic(ctx, r.RoundId, r.TopicId); err != nil {
		switch {
		case errors.Is(err, apperr.RoundNotFound):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundNotFound)
		case errors.Is(err, apperr.PackNotAuthor):
			return nil, twirp.PermissionDenied.Error(apperr.MsgPackNotAuthor)
		case errors.Is(err, apperr.RoundTopicNotAdded):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundTopicNotAdded)
		case errors.Is(err, apperr.RoundTopicAlreadyExists):
			return nil, twirp.InvalidArgument.Error(apperr.MsgRoundTopicAlreadyExists)
		}

		return nil, twirp.InternalError(err.Error())
	}

	return new(emptypb.Empty), nil
}
