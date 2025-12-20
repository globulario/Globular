package media

import (
	"encoding/json"
	"net/http"

	httplib "github.com/globulario/Globular/internal/gateway/http"
)

// SeasonEpisodeResolver abstracts how we obtain season/episode/series for a title.
type SeasonEpisodeResolver interface {
	ResolveSeasonEpisode(titleID string) (season int, episode int, seriesID string, err error)
}

// NewGetIMDBSeasonEpisode implements:
//
//	GET /get-imdb-season-episode?id=tt1234567
//
// Response:
//
//	200 {"season":<int>, "episode":<int>, "seriesId":"ttXXXXXXX"}
//	400 on errors
func NewGetIMDBSeasonEpisode(r SeasonEpisodeResolver) http.Handler {
	type resp struct {
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		SeriesID string `json:"seriesId"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		if id == "" {
			httplib.WriteJSONError(w, http.StatusBadRequest, "missing 'id' (IMDB title id like tt1234567)")
			return
		}
		season, episode, seriesID, err := r.ResolveSeasonEpisode(id)
		if err != nil {
			httplib.WriteJSONError(w, http.StatusBadRequest, "fail to resolve season/episode: "+err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp{
			Season:   season,
			Episode:  episode,
			SeriesID: seriesID,
		})
	})
}
