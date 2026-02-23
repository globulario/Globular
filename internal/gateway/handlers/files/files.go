package files

import "net/http"

type Deps struct {
	GetImages http.Handler
	Serve     http.Handler // /serve/*
	Upload    http.Handler // POST /file-upload
}

func Mount(mux *http.ServeMux, d Deps) {
	if d.GetImages != nil {
		mux.Handle("/api/get-images", d.GetImages)
	}
	if d.Serve != nil {
		mux.Handle("/api/serve/", http.StripPrefix("/api/serve/", d.Serve))
	}
	if d.Upload != nil {
		mux.Handle("/api/file-upload", d.Upload)
		mux.Handle("/uploads", d.Upload) // legacy alias used by globular-web-client
	}
}
