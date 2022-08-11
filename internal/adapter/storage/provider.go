package storage

import (
	"context"
	"io"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/answersuck/host/internal/config"
)

type provider struct {
	client      *minio.Client
	bucket      string
	storageHost string
	cdnHost     string
	useCDN      bool
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
		client:      c,
		bucket:      cfg.Bucket,
		cdnHost:     cfg.CDNHost,
		storageHost: cfg.Host,
		useCDN:      cfg.CDN,
	}, nil
}

func (p *provider) Upload(ctx context.Context, r io.Reader, name string, size int64, contentType string) (*url.URL, error) {
	opts := minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	res, err := p.client.PutObject(ctx, p.bucket, name, r, size, opts)
	if err != nil {
		return nil, err
	}

	return p.URL(res.Key), nil
}

func (p *provider) URL(filename string) *url.URL {
	if p.useCDN {
		return &url.URL{
			Scheme: "https",
			Host:   p.cdnHost,
			Path:   filename,
		}
	}

	return &url.URL{
		Scheme: "https",
		Host:   p.storageHost,
		Path:   p.bucket + "/" + filename,
	}
}
