package media

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httplib "github.com/globulario/Globular/internal/http"
)

// TitlesQuery captures common filters.
type TitlesQuery struct {
	Q      string
	Year   int    // optional
	Type   string // optional: movie|series|episode (provider decides)
	Limit  int
	Offset int
}

// TitlesSearcher abstracts how we fetch IMDB titles.
type TitlesSearcher interface {
	SearchIMDBTitles(q TitlesQuery) ([]map[string]any, error)
}

// NewGetIMDBTitles returns GET /get-imdb-titles.
// Response: 200 {"results":[ ... ]} or 400 on errors.
func NewGetIMDBTitles(s TitlesSearcher) http.Handler {
	type resp struct {
		Results []map[string]any `json:"results"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := TitlesQuery{
			Q:      strings.TrimSpace(r.URL.Query().Get("q")),
			Type:   strings.TrimSpace(r.URL.Query().Get("type")),
			Limit:  atoiDefault(r.URL.Query().Get("limit"), 20, 1, 200),
			Offset: atoiDefault(r.URL.Query().Get("offset"), 0, 0, 1_000_000),
		}
		if y := strings.TrimSpace(r.URL.Query().Get("year")); y != "" {
			if yy, err := strconv.Atoi(y); err == nil {
				q.Year = yy
			}
		}
		if q.Q == "" {
			httplib.WriteJSONError(w, http.StatusBadRequest, "missing 'q' query param")
			return
		}

		res, err := s.SearchIMDBTitles(q)
		if err != nil {
			httplib.WriteJSONError(w, http.StatusBadRequest, "fail to search IMDB titles: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		_ = enc.Encode(resp{Results: res})
	})
}

func atoiDefault(s string, def, min, max int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
