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

// WriteJSON sends any JSON payload using the specified status code.
// An empty payload is allowed but if data is nil and status is NoContent the body stays empty.
func WriteJSON(w stdhttp.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil || status == stdhttp.StatusNoContent {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}
