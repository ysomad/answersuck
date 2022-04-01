package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

type fileStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
}

func NewFileStorage(c *minio.Client, bucket, endpoint string) *fileStorage {
	return &fileStorage{
		client:   c,
		bucket:   bucket,
		endpoint: endpoint,
	}
}

func (s *fileStorage) Upload(ctx context.Context, f File) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  f.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	_, err := s.client.PutObject(ctx, s.bucket, f.Name, f.Reader, f.Size, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s", s.endpoint, s.bucket, f.Name), nil
}
