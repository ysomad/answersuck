package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ysomad/answersuck/internal/entity"
	pb "github.com/ysomad/answersuck/internal/gen/api/editor/v1"
	"github.com/ysomad/answersuck/internal/pkg/session"
	apptwirp "github.com/ysomad/answersuck/internal/twirp"
	"github.com/ysomad/answersuck/internal/twirp/common"
	"github.com/ysomad/answersuck/internal/twirp/hooks"
	"github.com/ysomad/answersuck/internal/twirp/middleware"
)

var (
	_ apptwirp.Handler = &TopicHandler{}
	_ pb.TopicService  = &TopicHandler{}
)

type TopicUseCase interface {
	Save(context.Context, entity.Topic) (topicID int32, err error)
}

type TopicHandler struct {
	topic   TopicUseCase
	session *session.Manager
}

func NewTopicHandler(uc TopicUseCase, sm *session.Manager) *TopicHandler {
	return &TopicHandler{
		topic:   uc,
		session: sm,
	}
}

func (h *TopicHandler) Handle(m *http.ServeMux) {
	s := pb.NewTopicServiceServer(h,
		twirp.WithServerHooks(hooks.WithSession(h.session)))
	m.Handle(s.PathPrefix(), middleware.WithSessionID(s))
}

func (h *TopicHandler) CreateTopic(
	ctx context.Context,
	r *pb.CreateTopicRequest) (*pb.CreateTopicResponse, error) {
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
