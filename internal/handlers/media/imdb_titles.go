package media

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
			http.Error(w, "missing 'q' query param", http.StatusBadRequest)
			return
		}

		res, err := s.SearchIMDBTitles(q)
		if err != nil {
			http.Error(w, "fail to search IMDB titles: "+err.Error(), http.StatusBadRequest)
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
