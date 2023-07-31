package round

import (
	"context"
	"fmt"
)

func (s *Service) RemoveTopic(ctx context.Context, roundID, topicID int32) error {
	if err := s.verifyPackAuthorship(ctx, roundID); err != nil {
		return fmt.Errorf("s.verifyPackAuthorship: %w", err)
	}

	if err := s.roundTopic.DeleteOne(ctx, roundID, topicID); err != nil {
		return fmt.Errorf("s.roundTopicDeleteOne: %w", err)
	}

	return nil
}
