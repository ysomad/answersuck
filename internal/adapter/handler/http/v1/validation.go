package v1

import (
	"errors"
	"io"
)

var errInvalidRequestBody = errors.New("invalid request body")

type ValidationModule interface {
	ValidateStruct(s any) error
	TranslateError(err error) map[string]string
	ValidateRequestBody(b io.ReadCloser, dest any) error
}
