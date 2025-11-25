package media

import (
	"encoding/json"
	"net/http"
)

// TrailerFetcher resolves an IMDb trailer video URL.
type TrailerFetcher interface {
	FetchIMDBTrailer(imdbID string) (pageURL string, imageURL string, videoSrc string, err error)
}

// NewGetIMDBTrailer returns GET /get-imdb-trailer?id=tt1234567
// Response: 200 {"id":"tt1234567","url":"https://..."} when found,
// 204 when no trailer is available, 400 on errors.
func NewGetIMDBTrailer(f TrailerFetcher) http.Handler {
	type resp struct {
		ID    string `json:"id"`
		URL   string `json:"url,omitempty"`
		Image string `json:"image,omitempty"`
		Video string `json:"video,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing 'id' (IMDb title id like tt1234567)", http.StatusBadRequest)
			return
		}
		url, image, video, err := f.FetchIMDBTrailer(id)
		if err != nil {
			http.Error(w, "fail to fetch IMDb trailer: "+err.Error(), http.StatusBadRequest)
			return
		}
		if url == "" && image == "" && video == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp{ID: id, URL: url, Image: image, Video: video})
	})
}
