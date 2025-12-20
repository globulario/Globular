package media

import (
	"net/http"

	httplib "github.com/globulario/Globular/internal/gateway/http"
)

// PosterFetcher can return the raw image OR a URL to redirect to.
type PosterFetcher interface {
	// FetchIMDBPoster returns (content, contentType, redirectURL, err).
	// Implementations can fill either content+contentType, or redirectURL.
	FetchIMDBPoster(imdbID, size string) (content []byte, contentType string, redirectURL string, err error)
}

// NewGetIMDBPoster returns GET /get-imdb-poster?id=tt1234567&size=...&mode=redirect|bytes
// - If provider returns content, we write it (200 + Content-Type).
// - Else if provider returns redirectURL, we redirect (307).
// - If mode=redirect is requested, we prefer redirect when available.
func NewGetIMDBPoster(p PosterFetcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			httplib.WriteJSONError(w, http.StatusBadRequest, "missing 'id' (IMDB id like tt1234567)")
			return
		}
		size := r.URL.Query().Get("size") // provider-specific (e.g. small, medium, large)
		mode := r.URL.Query().Get("mode") // optional: "redirect" or "bytes"

		content, ctype, url, err := p.FetchIMDBPoster(id, size)
		if err != nil {
			httplib.WriteJSONError(w, http.StatusBadRequest, "fail to get IMDB poster: "+err.Error())
			return
		}

		if mode == "redirect" && url != "" {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
		if len(content) > 0 && ctype != "" {
			w.Header().Set("Content-Type", ctype)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(content)
			return
		}
		if url != "" {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		// Nothing available
		w.WriteHeader(http.StatusNoContent)
	})
}
