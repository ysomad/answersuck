package media

import "io"

type (
	UploadDTO struct {
		Filename    string
		ContentType string
		Size        int64
		Buf         []byte
		AccountId   string
	}

	File struct {
		Reader      io.Reader
		Name        string
		Size        int64
		ContentType string
	}
)
