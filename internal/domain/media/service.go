package media

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/answersuck/host/internal/pkg/mime"
)

type (
	repository interface {
		Save(ctx context.Context, m Media) error
		FindById(ctx context.Context, mediaId string) (Media, error)
		FindMediaTypeById(ctx context.Context, mediaId string) (mime.Type, error)
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

func NewService(r repository, s fileStorage) *service {
	return &service{
		repo:    r,
		storage: s,
	}
}

func (s *service) UploadAndSave(ctx context.Context, m Media, size int64) (UploadedMediaDTO, error) {
	defer os.Remove(m.Filename)

	file, err := os.Open(m.Filename)
	if err != nil {
		return UploadedMediaDTO{}, fmt.Errorf("mediaService - UploadAndSave - os.Open: %w", err)
	}
	defer file.Close()

	url, err := s.storage.Upload(ctx, file, m.Filename, size, string(m.Type))
	if err != nil {
		return UploadedMediaDTO{}, fmt.Errorf("mediaService - UploadAndSave - s.storage.Upload: %w", err)
	}

	if err = s.repo.Save(ctx, m); err != nil {
		return UploadedMediaDTO{}, fmt.Errorf("mediaService - UploadAndSave - s.repo.Save: %w", err)
	}

	return UploadedMediaDTO{
		Id:        m.Id,
		URL:       url.String(),
		MediaType: m.Type,
	}, nil
}

func (s *service) GetMediaTypeById(ctx context.Context, mediaId string) (mime.Type, error) {
	t, err := s.repo.FindMediaTypeById(ctx, mediaId)
	if err != nil {
		return "", fmt.Errorf("mediaService - GetMimeTypeById - s.repo.FindMimeTypeById: %w", err)
	}

	return t, nil
}

func (s *service) GetById(ctx context.Context, mediaId string) (Media, error) {
	m, err := s.repo.FindById(ctx, mediaId)
	if err != nil {
		return Media{}, fmt.Errorf("mediaService - GetById - s.repo.FindById: %w", err)
	}

	return m, nil
}
