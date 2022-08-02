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
	b, _ := json.Marshal(apiErr{Message: err.Error()})
	w.WriteHeader(code)
	w.Write(b)
}

func writeValidationErr(w http.ResponseWriter, code int, err error, detail map[string]string) {
	b, _ := json.Marshal(validationErr{
		Message: err.Error(),
		Detail:  detail,
	})
	w.WriteHeader(code)
	w.Write(b)
}
