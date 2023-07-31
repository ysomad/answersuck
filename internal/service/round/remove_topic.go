package round

import (
	"context"
	"fmt"

	"github.com/ysomad/answersuck/internal/pkg/appctx"
	"github.com/ysomad/answersuck/internal/pkg/apperr"
)

func (s *Service) RemoveTopic(ctx context.Context, roundID, topicID int32) error {
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

	if err := s.roundTopic.DeleteOne(ctx, roundID, topicID); err != nil {
		return fmt.Errorf("s.roundTopicDeleteOne: %w", err)
	}

	return nil
}
