package answer

import "context"

type Repository interface {
	Create(ctx context.Context, a Answer) (Answer, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, a Answer) (Answer, error) {
	return s.repo.Create(ctx, a)
}
