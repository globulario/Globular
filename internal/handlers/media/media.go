package media

import "net/http"

// Deps lists the handlers to mount (all optional).
type Deps struct {
	GetIMDBTitles        http.Handler // /get-imdb-titles
	GetIMDBPoster        http.Handler // /get-imdb-poster
	GetIMDBSeasonEpisode http.Handler // /get-imdb-season-episode   <-- NEW
	GetIMDBTrailer       http.Handler // /get-imdb-trailer
}

// Mount registers the provided handlers.
func Mount(mux *http.ServeMux, d Deps) {
	if d.GetIMDBTitles != nil {
		mux.Handle("/api/get-imdb-titles", d.GetIMDBTitles)
	}
	if d.GetIMDBPoster != nil {
		mux.Handle("/api/get-imdb-poster", d.GetIMDBPoster)
	}
	if d.GetIMDBSeasonEpisode != nil {
		mux.Handle("/api/get-imdb-season-episode", d.GetIMDBSeasonEpisode)
	}
	if d.GetIMDBTrailer != nil {
		mux.Handle("/api/get-imdb-trailer", d.GetIMDBTrailer)
	}
}
