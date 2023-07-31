package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (s *Service) AddTopic(ctx context.Context, roundID, topicID int32) error {
	if err := s.verifyPackAuthorship(ctx, roundID); err != nil {
		return fmt.Errorf("s.verifyPackAuthorship: %w", err)
	}

	roundTopics, err := s.roundTopic.GetAll(ctx, roundID)
	if err != nil {
		return fmt.Errorf("s.roundTopic.GetAll: %w", err)
	}

	if len(roundTopics) == entity.MaxRoundTopics {
		return apperr.RoundTopicNotAdded
	}

	if err := s.roundTopic.Save(ctx, roundID, topicID); err != nil {
		return fmt.Errorf("s.roundTopic.Save: %w", err)
	}

	return nil
}
