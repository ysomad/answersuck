package v1

import (
	"encoding/json"
	"errors"
	"net/http"
)

var errInvalidRequestBody = errors.New("invalid request body")

type validationErr struct {
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail"`
}

type apiErr struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func writeErr(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiErr{Message: err.Error()})
}

func writeValidationErr(w http.ResponseWriter, code int, err error, detail map[string]string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(validationErr{
		Message: err.Error(),
		Detail:  detail,
	})
}
