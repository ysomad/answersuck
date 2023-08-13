package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (s *Service) AddTopic(ctx context.Context, roundID, topicID int32) (int32, error) {
	if err := s.pack.VerifyRoundAuthorship(ctx, roundID); err != nil {
		return 0, fmt.Errorf("error verifying round authorship: %w", err)
	}

	roundTopics, err := s.roundTopic.GetAll(ctx, roundID)
	if err != nil {
		return 0, fmt.Errorf("error getting round topics: %w", err)
	}

	if len(roundTopics) == entity.MaxRoundTopics {
		return 0, apperr.RoundTopicNotAdded
	}

	roundTopicID, err := s.roundTopic.Save(ctx, roundID, topicID)
	if err != nil {
		return 0, fmt.Errorf("error saving round topic: %w", err)
	}

	return roundTopicID, nil
}
