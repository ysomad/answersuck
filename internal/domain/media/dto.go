package media

import "io"

type (
	File struct {
		Reader      io.Reader
		Name        string
		Size        int64
		ContentType string
	}

	WithURL struct {
		Id        string `json:"id"`
		URL       string `json:"url"`
		MediaType Type   `json:"type"`
	}
)
