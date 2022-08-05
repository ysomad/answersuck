package v1

import (
	"encoding/json"
	"net/http"
)

type listResponse struct {
	Result any `json:"result"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
	w.WriteHeader(code)
}
