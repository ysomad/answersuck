package repository

import (
	"context"

	"github.com/answersuck/vault/internal/domain"
	"github.com/answersuck/vault/pkg/postgres"
)

type sessionRepository struct {
	*postgres.Client
}

func NewSessionRepository(pg *postgres.Client) *sessionRepository {
	return &sessionRepository{pg}
}

func (r *sessionRepository) Create(ctx context.Context, s *domain.Session) error {
	panic("implement")

	return nil
}

func (r *sessionRepository) FindById(ctx context.Context, sid string) (*domain.Session, error) {
	panic("implement")

	return nil, nil
}

func (r *sessionRepository) FindAll(ctx context.Context, aid string) ([]*domain.Session, error) {
	panic("implement")

	return nil, nil
}

func (r *sessionRepository) Delete(ctx context.Context, sid string) error {
	panic("implement")

	return nil
}

func (r *sessionRepository) DeleteAll(ctx context.Context, aid, sid string) error {
	panic("implement")

	return nil
}
