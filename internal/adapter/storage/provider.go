package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/answersuck/host/internal/config"
	"github.com/answersuck/host/internal/domain/media"
)

type provider struct {
	client        *minio.Client
	bucket        string
	storageDomain string
	cdnDomain     string
	useCDN        bool
}

func NewProvider(cfg *config.FileStorage) (*provider, error) {
	c, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		return nil, err
	}

	return &provider{
		client:        c,
		bucket:        cfg.Bucket,
		cdnDomain:     cfg.CDNDomain,
		storageDomain: cfg.Domain,
		useCDN:        cfg.CDN,
	}, nil
}

func (p *provider) Upload(ctx context.Context, f media.File) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  f.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	res, err := p.client.PutObject(ctx, p.bucket, f.Name, f.Reader, f.Size, opts)
	if err != nil {
		return "", err
	}

	if p.useCDN {
		return p.cdnURL(p.cdnDomain, res.Key), nil
	}

	return p.storageURL(p.storageDomain, p.bucket, res.Key), nil
}

// Private

func (p *provider) cdnURL(domain, filename string) string {
	return fmt.Sprintf("https://%s/%s", domain, filename)
}

func (p *provider) storageURL(domain, bucket, filename string) string {
	return fmt.Sprintf("https://%s/%s/%s", domain, bucket, filename)
}
