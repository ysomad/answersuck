package media

import (
	"context"
	"fmt"
	"net/url"
	"os"
)

type (
	repository interface {
		Save(ctx context.Context, m Media) error
		FindMediaTypeById(ctx context.Context, mediaId string) (string, error)
	}

	fileStorage interface {
		Upload(ctx context.Context, f File) (url url.URL, err error)
		URL(filename string) url.URL
	}
)

type service struct {
	repo    repository
	storage fileStorage
}

func NewService(r repository, s fileStorage) *service {
	return &service{
		repo:    r,
		storage: s,
	}
}

func (s *service) UploadAndSave(ctx context.Context, m Media, size int64) (WithURL, error) {
	defer m.removeTmpFile()

	f, err := os.Open(m.Filename)
	if err != nil {
		return WithURL{}, fmt.Errorf("mediaService - UploadAndSave - os.Open: %w", err)
	}
	defer f.Close()

	url, err := s.storage.Upload(ctx, File{
		Reader:      f,
		Name:        m.Filename,
		Size:        size,
		ContentType: string(m.Type),
	})
	if err != nil {
		return WithURL{}, fmt.Errorf("mediaService - UploadAndSave - s.storage.Upload: %w", err)
	}

	if err = s.repo.Save(ctx, m); err != nil {
		return WithURL{}, fmt.Errorf("mediaService - UploadAndSave - s.repo.Save: %w", err)
	}

	return WithURL{
		Id:        m.Id,
		URL:       url.String(),
		MediaType: m.Type,
	}, nil
}

func (s *service) GetMediaTypeById(ctx context.Context, mediaId string) (string, error) {
	t, err := s.repo.FindMediaTypeById(ctx, mediaId)
	if err != nil {
		return "", fmt.Errorf("mediaService - GetMimeTypeById - s.repo.FindMimeTypeById: %w", err)
	}

	return t, nil
}
