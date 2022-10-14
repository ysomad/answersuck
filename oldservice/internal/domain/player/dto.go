package player

import "github.com/ysomad/answersuck-backend/internal/pkg/mime"

type UploadAvatarDTO struct {
	AccountId   string
	Filename    string
	FileSize    int64
	ContentType mime.Type
}
