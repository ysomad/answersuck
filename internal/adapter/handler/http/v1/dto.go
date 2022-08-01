package v1

import (
	"encoding/json"
	"net/http"
)

type listResp struct {
	Result any `json:"result"`
}

func writeJSON(w http.ResponseWriter, v any) {
	b, _ := json.Marshal(v)
	w.Write(b)
}

func writeList(w http.ResponseWriter, v any) {
	b, _ := json.Marshal(listResp{Result: v})
	w.Write(b)
}
