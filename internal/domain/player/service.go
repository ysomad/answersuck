package player

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"
)

type (
	repository interface {
		FindByNickname(ctx context.Context, nickname string) (Player, error)
		SetAvatar(ctx context.Context, a Avatar) error
	}

	fileStorage interface {
		Upload(ctx context.Context, r io.Reader, name string, size int64, contentType string) (*url.URL, error)
		URL(filename string) *url.URL
	}
)

type service struct {
	repo    repository
	storage fileStorage
}

func NewService(r repository, p fileStorage) *service {
	return &service{
		repo:    r,
		storage: p,
	}
}

func (s *service) GetByNickname(ctx context.Context, nickname string) (Detailed, error) {
	p, err := s.repo.FindByNickname(ctx, nickname)
	if err != nil {
		return Detailed{}, fmt.Errorf("playerService - GetByNickname - s.repo.FindByNickname: %w", err)
	}

	return NewDetailed(p, s.storage), nil
}

func (s *service) UploadAvatar(ctx context.Context, dto UploadAvatarDTO) error {
	defer os.Remove(dto.Filename)

	file, err := os.Open(dto.Filename)
	if err != nil {
		return fmt.Errorf("playerService - UploadAvatar - os.Open: %w", err)
	}
	defer file.Close()

	if _, err = s.storage.Upload(ctx, file, dto.Filename, dto.FileSize, string(dto.ContentType)); err != nil {
		return fmt.Errorf("playerService - UploadAvatar - s.storage.Upload: %w", err)
	}

	if err := s.repo.SetAvatar(ctx, Avatar{
		AccountId: dto.AccountId,
		Filename:  dto.Filename,
		UpdatedAt: time.Now(),
	}); err != nil {
		return fmt.Errorf("playerService - UploadAvatar - s.repo.SetAvatar: %w", err)
	}

	return nil
}
