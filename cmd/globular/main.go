package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	httplib "Globular/internal/http"
	"Globular/internal/server"
)

func main() {
	// Flags (minimal for now — we’ll wire more later)
	var (
		httpAddr  = flag.String("http", ":8080", "HTTP listen address (empty to disable)")
		httpsAddr = flag.String("https", "", "HTTPS listen address (empty to disable)")
	)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Build router (middleware, metrics, health)
	mux := httplib.NewRouter(logger, httplib.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		RateRPS:        20,
		RateBurst:      80,
	})

	// ---- Mount existing handlers here (no behavior change) ----
	// mux.Handle("/getConfig", http.HandlerFunc(getConfigHanldler))
	// mux.Handle("/get-images", http.HandlerFunc(GetImagesHandler))
	// mux.Handle("/file-upload", http.HandlerFunc(FileUploadHandler))
	// mux.Handle("/serve", http.HandlerFunc(ServeFileHandler))
	// mux.Handle("/get-imdb-titles", http.HandlerFunc(getImdbTitlesHanldler))
	// -----------------------------------------------------------

	// Server supervisor (HTTP/HTTPS + graceful shutdown)
	sup := server.Supervisor{
		Logger:            logger,
		HTTPAddr:          *httpAddr,
		HTTPSAddr:         *httpsAddr,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if err := sup.Start(mux); err != nil {
		logger.Error("start failed", "err", err)
		os.Exit(1)
	}

	// Graceful stop
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	_ = sup.Stop(context.Background())
}
