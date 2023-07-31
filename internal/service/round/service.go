package round

import (
	"context"

	"github.com/ysomad/answersuck/internal/entity"
)

type packService interface {
	GetOne(ctx context.Context, packID int32) (*entity.Pack, error)
}

type roundTopicService interface {
	GetAll(ctx context.Context, roundID int32) ([]entity.Topic, error)
	Save(ctx context.Context, roundID, topicID int32) error
	DeleteOne(ctx context.Context, roundID, topicID int32) error
}

type repository interface {
	Save(ctx context.Context, round entity.Round) (int32, error)
	UpdateOne(context.Context, entity.Round) error
	GetPackAuthor(ctx context.Context, roundID int32) (string, error)
}

type Service struct {
	repo       repository
	pack       packService
	roundTopic roundTopicService
}

func NewService(r repository, ps packService, rts roundTopicService) *Service {
	return &Service{
		repo:       r,
		pack:       ps,
		roundTopic: rts,
	}
}
