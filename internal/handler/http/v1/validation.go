package v1

import "errors"

var errInvalidRequestBody = errors.New("invalid request body")

type ValidationModule interface {
	ValidateStruct(s any) error
	TranslateError(err error) map[string]string
}
