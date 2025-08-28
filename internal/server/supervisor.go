package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Supervisor struct {
	Logger            *slog.Logger
	HTTPAddr          string
	HTTPSAddr         string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	httpSrv  *http.Server
	httpsSrv *http.Server
}

func (s *Supervisor) Start(handler http.Handler) error {
	if s.HTTPAddr != "" {
		s.httpSrv = &http.Server{
			Addr:              s.HTTPAddr,
			Handler:           handler,
			ReadHeaderTimeout: s.ReadHeaderTimeout,
			ReadTimeout:       s.ReadTimeout,
			WriteTimeout:      s.WriteTimeout,
			IdleTimeout:       s.IdleTimeout,
		}
		go func() {
			s.Logger.Info("http listen", "addr", s.HTTPAddr)
			if err := s.httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
				s.Logger.Error("http server error", "err", err)
			}
		}()
	}
	// HTTPS wiring will come later (TLS or ACME)
	return nil
}

func (s *Supervisor) Stop(ctx context.Context) error {
	if s.httpsSrv != nil {
		_ = s.httpsSrv.Shutdown(ctx)
	}
	if s.httpSrv != nil {
		_ = s.httpSrv.Shutdown(ctx)
	}
	return nil
}
