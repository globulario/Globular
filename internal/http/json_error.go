package http

import (
	"encoding/json"
	stdhttp "net/http"
)

type jsonErrorResp struct {
	Error string `json:"error"`
}

// WriteJSONError writes a structured JSON error response while preserving the
// provided HTTP status code. For StatusNoContent we only set the status so
// the response stays empty as required by the spec.
func WriteJSONError(w stdhttp.ResponseWriter, status int, message string) {
	if message == "" {
		message = stdhttp.StatusText(status)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == stdhttp.StatusNoContent {
		return
	}
	_ = json.NewEncoder(w).Encode(jsonErrorResp{Error: message})
}
