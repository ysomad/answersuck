package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/answersuck/vault/internal/domain"
)

type SessionRepository struct {
	*redis.Client
}

func NewSessionRepository(r *redis.Client) *SessionRepository {
	return &SessionRepository{r}
}

func (r *SessionRepository) Create(ctx context.Context, s *domain.Session) error {
	b, err := s.MarshalBinary()
	if err != nil {
		return err
	}

	if err = r.Set(ctx, s.Id, b, s.Expiration).Err(); err != nil {
		return fmt.Errorf("r.Set: %w", err)
	}

	return nil
}

func (r *SessionRepository) FindById(ctx context.Context, sid string) (*domain.Session, error) {
	var s domain.Session

	if err := r.Get(ctx, sid).Scan(&s); err != nil {
		return nil, fmt.Errorf("r.Get.Scan: %w", err)
	}

	return &s, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sid string) error {
	if err := r.Del(ctx, sid).Err(); err != nil {
		return fmt.Errorf("r.Del: %w", err)
	}

	return nil
}
