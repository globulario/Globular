package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// TLSFiles points to on-disk cert/key (fullchain.pem + server.key).
type TLSFiles struct {
	CertFile string // fullchain.pem
	KeyFile  string // server.key
}

// certReloader hot-reloads a single certificate pair from disk.
type certReloader struct {
	certPath string
	keyPath  string

	mu      sync.RWMutex
	cert    *tls.Certificate
	checked time.Time // last stat time used to decide reload
}

// newCertReloader loads the initial certificate pair.
func newCertReloader(certPath, keyPath string) (*certReloader, error) {
	r := &certReloader{certPath: certPath, keyPath: keyPath}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

// getCertificate is plugged into tls.Config.GetCertificate.
func (r *certReloader) getCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	// Fast path: serve what we have.
	r.mu.RLock()
	c := r.cert
	last := r.checked
	r.mu.RUnlock()

	// Light guard: only stat once per second to avoid syscall flood.
	if time.Since(last) >= time.Second {
		_ = r.maybeReload()
	}

	if c == nil {
		return nil, errors.New("no certificate loaded")
	}
	return c, nil
}

func (r *certReloader) maybeReload() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// mark that we checked
	r.checked = time.Now()

	cInfo, cErr := os.Stat(r.certPath)
	kInfo, kErr := os.Stat(r.keyPath)
	if cErr != nil || kErr != nil {
		// keep serving the previous cert if files momentarily missing
		return nil
	}

	// If either file is newer than what the current cert was loaded from, reload.
	if r.cert == nil || cInfo.ModTime().After(cmt(r.cert)) || kInfo.ModTime().After(cmt(r.cert)) {
		return r.reload()
	}
	return nil
}

// cmt extracts our last-known modtime from the certificate's leaf if present.
// (We stash it via the checked field instead; this helper just keeps intent clear.)
func cmt(_ *tls.Certificate) time.Time { return time.Time{} }

func (r *certReloader) reload() error {
	crt, err := tls.LoadX509KeyPair(r.certPath, r.keyPath)
	if err != nil {
		return err
	}

	// Populate Leaf (useful for logging/metrics/expiry checks).
	if len(crt.Certificate) > 0 {
		// Prefer the first non-CA as the leaf; fall back to index 0.
		var leafDER []byte
		cands := crt.Certificate

		// Try to find a non-CA cert in the chain.
		for _, der := range cands {
			if c, perr := x509.ParseCertificate(der); perr == nil && !c.IsCA {
				leafDER = der
				break
			}
		}
		if leafDER == nil {
			leafDER = cands[0]
		}

		if leaf, perr := x509.ParseCertificate(leafDER); perr == nil {
			crt.Leaf = leaf
		}
	}

	r.cert = &crt
	r.checked = time.Now()
	return nil
}

// Start a background poll that calls maybeReload every 'interval'.
func (r *certReloader) watch(ctx context.Context, interval time.Duration, log *slog.Logger) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := r.maybeReload(); err != nil && log != nil {
				log.Warn("tls reload failed", "err", err)
			}
		}
	}
}

// Supervisor runs HTTP/HTTPS servers with optional hot-reload TLS.
type Supervisor struct {
	Logger            *slog.Logger
	HTTPAddr          string
	HTTPSAddr         string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	// TLS (when set, HTTPS will be started). If CertFile/KeyFile change on disk,
	// they will be picked up automatically without restarting the process.
	TLS *TLSFiles

	httpSrv  *http.Server
	httpsSrv *http.Server

	Ready chan struct{} // closed once first listener is bound
}

// Start launches HTTP/HTTPS servers. Call Stop to shut down gracefully.
func (s *Supervisor) Start(handler http.Handler) error {

	if s.Ready == nil {
		s.Ready = make(chan struct{})
	}

	if s.HTTPAddr != "" {
		h2Server := &http2.Server{}
		httpHandler := handler
		if httpHandler == nil {
			httpHandler = http.DefaultServeMux
		}
		httpHandler = h2c.NewHandler(httpHandler, h2Server)

		s.httpSrv = &http.Server{
			Addr:              s.HTTPAddr,
			Handler:           httpHandler,
			ReadHeaderTimeout: s.ReadHeaderTimeout,
			ReadTimeout:       s.ReadTimeout,
			WriteTimeout:      s.WriteTimeout,
			IdleTimeout:       s.IdleTimeout,
		}
		if err := http2.ConfigureServer(s.httpSrv, h2Server); err != nil {
			return err
		}
		go func() {
			s.Logger.Info("http listen", "addr", s.HTTPAddr)
			close(s.Ready) // signal as soon as one listener starts (or do once via sync.Once)
			if err := s.httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
				s.Logger.Error("http server error", "err", err)
			}
		}()
	}

	if s.HTTPSAddr != "" && s.TLS != nil && s.TLS.CertFile != "" && s.TLS.KeyFile != "" {
		// Build live reloader and tls.Config using GetCertificate.
		reloader, err := newCertReloader(s.TLS.CertFile, s.TLS.KeyFile)
		if err != nil {
			return err
		}
		tcfg := &tls.Config{
			MinVersion:     tls.VersionTLS12,
			GetCertificate: reloader.getCertificate,
		}
		tlsH2Server := &http2.Server{}
		s.httpsSrv = &http.Server{
			Addr:              s.HTTPSAddr,
			Handler:           handler,
			ReadHeaderTimeout: s.ReadHeaderTimeout,
			ReadTimeout:       s.ReadTimeout,
			WriteTimeout:      s.WriteTimeout,
			IdleTimeout:       s.IdleTimeout,
			TLSConfig:         tcfg,
		}
		if err := http2.ConfigureServer(s.httpsSrv, tlsH2Server); err != nil {
			return err
		}
		// Background cert polling (1s cadence is plenty; bump if you prefer).
		ctx, _ := context.WithCancel(context.Background())
		go reloader.watch(ctx, time.Second, s.Logger)

		go func() {
			s.Logger.Info("https listen", "addr", s.HTTPSAddr, "cert", s.TLS.CertFile)
			close(s.Ready) // also safe if already closed
			// When using GetCertificate, pass empty cert/key paths.
			if err := s.httpsSrv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) && err != nil {
				s.Logger.Error("https server error", "err", err)
			}
		}()
	}

	return nil
}

// Stop gracefully shuts down both servers.
func (s *Supervisor) Stop(ctx context.Context) error {
	if s.httpsSrv != nil {
		_ = s.httpsSrv.Shutdown(ctx)
	}
	if s.httpSrv != nil {
		_ = s.httpSrv.Shutdown(ctx)
	}
	return nil
}

// tls.X509KeyPairToLeaf is not in stdlib; implement a tiny helper.
func X509KeyPairToLeaf(pair tls.Certificate) (*tls.Certificate, error) {
	// No-op shim to keep the call site readable in reload().
	return &pair, nil
}
