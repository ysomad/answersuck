package packages

import (
	"context"
	"fmt"
	"time"
)

type repository interface {
	Save(ctx context.Context, p Package) (packageId uint32, err error)
}

type service struct {
	repo repository
}

func NewService(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(ctx context.Context, p CreateParams) (uint32, error) {
	now := time.Now()
	packageId, err := s.repo.Save(ctx, Package{
		Name:        p.Name,
		Description: p.Description,
		Published:   false,
		LanguageId:  p.LanguageId,
		Tags:        p.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return 0, fmt.Errorf("packageService - Create - s.repo.Save: %w", err)
	}

	return packageId, nil
}
