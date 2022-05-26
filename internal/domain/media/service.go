package media

import (
	"context"
	"fmt"
	"os"
	"time"
)

type (
	Repository interface {
		Save(ctx context.Context, m Media) (Media, error)
		FindMimeTypeById(ctx context.Context, mediaId string) (string, error)
	}

	Storage interface {
		Upload(ctx context.Context, f File) (string, error)
	}
)

type service struct {
	repo    Repository
	storage Storage
}

func NewService(r Repository, s Storage) *service {
	return &service{
		repo:    r,
		storage: s,
	}
}

func (s *service) UploadAndSave(ctx context.Context, dto *UploadDTO) (Media, error) {
	t := MimeType(dto.ContentType)

	if !t.valid() {
		return Media{}, fmt.Errorf("mediaService - UploadAndSave: %w", ErrInvalidMimeType)
	}

	m := Media{
		AccountId: dto.AccountId,
		Type:      t,
		CreatedAt: time.Now(),
	}

	if err := m.generateId(); err != nil {
		return Media{}, fmt.Errorf("mediaService - UploadAndSave - m.GenerateId: %w", err)
	}

	filename := m.filenameFromId(dto.Filename)

	tmp, err := m.newTempFile(filename, dto.Buf)
	if err != nil {
		return Media{}, fmt.Errorf("mediaService - UploadAndSave - m.NewFileFromBuffer: %w", err)
	}

	defer m.deleteTempFile(filename)
	defer tmp.Close()

	f, err := os.Open(filename)
	if err != nil {
		return Media{}, fmt.Errorf("mediaService - UploadAndSave - os.Open: %w", err)
	}

	defer f.Close()

	url, err := s.storage.Upload(ctx, File{
		Reader:      f,
		Name:        filename,
		Size:        dto.Size,
		ContentType: dto.ContentType,
	})

	if err != nil {
		return Media{}, fmt.Errorf("mediaService - UploadAndSave - s.storage.Upload: %w", err)
	}

	m.URL = url

	m, err = s.repo.Save(ctx, m)
	if err != nil {
		return Media{}, fmt.Errorf("mediaService - UploadAndSave - s.repo.Save: %w", err)
	}

	return m, nil
}

func (s *service) GetMimeTypeById(ctx context.Context, mediaId string) (string, error) {
	t, err := s.repo.FindMimeTypeById(ctx, mediaId)
	if err != nil {
		return "", fmt.Errorf("mediaService - GetMimeTypeById - s.repo.FindMimeTypeById: %w", err)
	}

	return t, nil
}
