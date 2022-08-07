package answer

import (
	"errors"

	"github.com/answersuck/host/internal/domain/media"
)

var (
	ErrMediaTypeNotAllowed = errors.New("not allowed media type for answer")
	ErrMediaNotFound       = errors.New("media with provided id not found")
)

type Answer struct {
	Id      int     `json:"id"`
	Text    string  `json:"text"`
	MediaId *string `json:"mediaId"`
}

var allowedMediaType = [3]media.Type{media.TypeImageJPEG, media.TypeImagePNG, media.TypeImageWEBP}

// mediaTypeAllowed checks if media for answer is in allowed media type array
func mediaTypeAllowed(mt string) bool {
	var allowed bool
	for _, t := range allowedMediaType {
		if media.Type(mt) == t {
			allowed = true
			break
		}
	}
	return allowed
}
