package player

import (
	"github.com/answersuck/host/internal/pkg/mime"
)

type UploadAvatarDTO struct {
	AccountId   string
	Filename    string
	FileSize    int64
	ContentType mime.Type
}
