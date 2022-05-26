package answer

import (
	"errors"
)

var (
	ErrMimeTypeNotAllowed = errors.New("not allowed mime type for answer")
	ErrMediaNotFound      = errors.New("media with provided id not found")
)
