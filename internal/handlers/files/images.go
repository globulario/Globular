package files

import (
	"encoding/json"
	"net/http"
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
			http.Error(w, "missing 'dir' query param", http.StatusBadRequest)
			return
		}
		files, err := l.ListImages(dir)
		if err != nil {
			http.Error(w, "fail to list images: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp{Files: files})
	})
}
