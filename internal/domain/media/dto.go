package media

import (
	"io"

	"github.com/answersuck/host/internal/pkg/mime"
)

type (
	File struct {
		Reader      io.Reader
		Name        string
		Size        int64
		ContentType string
	}

	UploadedMediaDTO struct {
		Id        string    `json:"id"`
		URL       string    `json:"url"`
		MediaType mime.Type `json:"type"`
	}
)
