package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/entity"
	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (s *Service) AddTopic(ctx context.Context, roundID, topicID int32) error {
	nickname, ok := appctx.GetNickname(ctx)
	if !ok {
		return apperr.Unauthorized
	}

	packAuthor, err := s.repo.GetPackAuthor(ctx, roundID)
	if err != nil {
		return fmt.Errorf("s.repo.GetPackAuthor: %w", err)
	}

	if packAuthor != nickname {
		return apperr.PackNotAuthor
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
