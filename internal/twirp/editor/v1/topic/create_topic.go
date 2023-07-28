package topic

import (
	"context"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateTopic(ctx context.Context, r *pb.CreateTopicRequest) (*pb.CreateTopicResponse, error) {
	session, err := common.CheckPlayerVerification(ctx)
	if err != nil {
		return nil, err
	}

	if r.TopicTitle == "" {
		return nil, twirp.RequiredArgumentError("topic_title")
	}

	if err := r.Validate(); err != nil {
		return nil, twirp.InvalidArgument.Error(err.Error())
	}

	topic := entity.Topic{
		Title:      r.TopicTitle,
		Author:     session.User.ID,
		CreateTime: time.Now(),
	}

	topic.ID, err = h.topic.Save(ctx, topic)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	return &pb.CreateTopicResponse{
		Topic: &pb.Topic{
			Id:         topic.ID,
			Title:      topic.Title,
			CreateTime: timestamppb.New(topic.CreateTime),
		},
	}, nil
}
