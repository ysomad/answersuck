package player

import (
	"context"
	"fmt"
	"net/url"
)

type (
	repository interface {
		FindByNickname(ctx context.Context, nickname string) (Player, error)
	}

	mediaProvider interface {
		URL(filename string) *url.URL
	}
)

type service struct {
	repo  repository
	media mediaProvider
}

func NewService(r repository, p mediaProvider) *service {
	return &service{
		repo:  r,
		media: p,
	}
}

func (s *service) GetByNickname(ctx context.Context, nickname string) (Player, error) {
	p, err := s.repo.FindByNickname(ctx, nickname)
	if err != nil {
		return Player{}, fmt.Errorf("playerService - GetByNickname - s.repo.FindByNickname: %w", err)
	}

	p.setAvatarURL(s.media)

	return p, nil
}
