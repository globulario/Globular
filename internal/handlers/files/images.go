package files

import (
	"encoding/json"
	"net/http"

	httplib "github.com/globulario/Globular/internal/http"
)

type ImageLister interface {
	ListImages(dir string) ([]string, error)
}

func NewGetImages(l ImageLister) http.Handler {
	type resp struct {
		Files []string `json:"files"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("dir")
		if dir == "" {
			httplib.WriteJSONError(w, http.StatusBadRequest, "missing 'dir' query param")
			return
		}
		files, err := l.ListImages(dir)
		if err != nil {
			httplib.WriteJSONError(w, http.StatusBadRequest, "fail to list images: "+err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp{Files: files})
	})
}
