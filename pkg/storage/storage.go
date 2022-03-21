package storage

import (
	"context"
	"io"
)

type File struct {
	Reader      io.Reader
	Name        string
	Size        int64
	ContentType string
}

type Uploader interface {
	// Upload uploads provided file and returns URL to it or error
	Upload(ctx context.Context, f File) (string, error)
}
