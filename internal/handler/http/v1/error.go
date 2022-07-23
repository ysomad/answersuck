package v1

import (
	"encoding/json"
	"net/http"
)

type detailedErr[D string | map[string]string] struct {
	Message string `json:"message"`
	Details D      `json:"details"`
}

type msgResp struct {
	Message string `json:"message"`
}

type detailedMsgResp struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func writeError(w http.ResponseWriter, code int, err error) {
	b, _ := json.Marshal(msgResp{Message: err.Error()})
	w.WriteHeader(code)
	w.Write(b)
}

func writeDetailedError(w http.ResponseWriter, code int, err error, details map[string]string) {
	b, _ := json.Marshal(detailedMsgResp{
		Message: err.Error(),
		Details: details,
	})
	w.WriteHeader(code)
	w.Write(b)
}
