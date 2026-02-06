package watchers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// CertReconciler watches for certificate changes via etcd and filesystem events.
// Replaces polling with event-driven reconciliation (v1 conformance INV-6).
type CertReconciler struct {
	logger     *slog.Logger
	etcdClient *clientv3.Client

	// Event channels
	certChanged chan struct{} // Signal channel for certificate changes
	acmeChanged chan struct{} // Signal channel for ACME certificate changes

	// State tracking
	mu                 sync.RWMutex
	lastCertGeneration uint64
	lastACMECertHash   string
	lastACMEKeyHash    string

	// Configuration
	domain       string // For etcd key (TODO: replace with cert ID in future)
	acmeCertPath string
	acmeKeyPath  string

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewCertReconciler creates a certificate reconciler with event-driven watching.
func NewCertReconciler(logger *slog.Logger, etcdClient *clientv3.Client, domain, acmeCertPath, acmeKeyPath string) *CertReconciler {
	if logger == nil {
		logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CertReconciler{
		logger:       logger,
		etcdClient:   etcdClient,
		certChanged:  make(chan struct{}, 1),
		acmeChanged:  make(chan struct{}, 1),
		domain:       domain,
		acmeCertPath: acmeCertPath,
		acmeKeyPath:  acmeKeyPath,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins watching for certificate changes.
func (r *CertReconciler) Start() error {
	// Start etcd watcher for internal certificates
	if r.etcdClient != nil && r.domain != "" {
		r.wg.Add(1)
		go r.watchEtcdCertificate()
	}

	// Start filesystem watcher for ACME certificates
	if r.acmeCertPath != "" && r.acmeKeyPath != "" {
		r.wg.Add(1)
		go r.watchACMECertificates()
	}

	return nil
}

// Stop gracefully shuts down the reconciler.
func (r *CertReconciler) Stop() {
	r.cancel()
	r.wg.Wait()
}

// CertChangedChan returns a channel that signals internal certificate changes.
func (r *CertReconciler) CertChangedChan() <-chan struct{} {
	return r.certChanged
}

// ACMEChangedChan returns a channel that signals ACME certificate changes.
func (r *CertReconciler) ACMEChangedChan() <-chan struct{} {
	return r.acmeChanged
}

// watchEtcdCertificate watches etcd for certificate generation changes.
// v1 Conformance (INV-6.1): Event-driven certificate rotation via etcd Watch.
func (r *CertReconciler) watchEtcdCertificate() {
	defer r.wg.Done()

	// VIOLATION INV-1.8: Domain-based persistent state key
	// TODO: Use /globular/pki/certs/{cert_id} instead
	key := fmt.Sprintf("/globular/pki/bundles/%s", r.domain)

	r.logger.Info("starting etcd certificate watcher", "key", key)

	// Initialize current generation
	if err := r.initializeCertGeneration(key); err != nil {
		r.logger.Warn("failed to initialize certificate generation", "err", err)
	}

	// Create watch channel
	watchChan := r.etcdClient.Watch(r.ctx, key)

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("etcd certificate watcher stopped")
			return

		case watchResp := <-watchChan:
			if watchResp.Err() != nil {
				r.logger.Error("etcd watch error", "err", watchResp.Err())
				// Retry with backoff
				time.Sleep(5 * time.Second)
				watchChan = r.etcdClient.Watch(r.ctx, key)
				continue
			}

			for _, event := range watchResp.Events {
				if event.Type == clientv3.EventTypePut {
					if r.handleCertGenerationChange(event.Kv.Value) {
						// Signal certificate change (non-blocking)
						select {
						case r.certChanged <- struct{}{}:
						default:
							// Already pending, skip
						}
					}
				}
			}
		}
	}
}

// initializeCertGeneration reads the current certificate generation from etcd.
func (r *CertReconciler) initializeCertGeneration(key string) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	resp, err := r.etcdClient.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		r.logger.Info("no certificate bundle found in etcd", "key", key)
		return nil
	}

	var payload struct {
		Generation uint64 `json:"generation"`
	}
	if err := json.Unmarshal(resp.Kvs[0].Value, &payload); err != nil {
		return fmt.Errorf("parse certificate generation: %w", err)
	}

	r.mu.Lock()
	r.lastCertGeneration = payload.Generation
	r.mu.Unlock()

	r.logger.Info("certificate generation initialized", "generation", payload.Generation)
	return nil
}

// handleCertGenerationChange processes a certificate generation change event.
// Returns true if generation actually changed.
func (r *CertReconciler) handleCertGenerationChange(value []byte) bool {
	var payload struct {
		Generation uint64 `json:"generation"`
	}
	if err := json.Unmarshal(value, &payload); err != nil {
		r.logger.Warn("failed to parse certificate generation", "err", err)
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if payload.Generation != r.lastCertGeneration {
		r.logger.Info("certificate generation changed",
			"old", r.lastCertGeneration,
			"new", payload.Generation,
			"domain", r.domain)
		r.lastCertGeneration = payload.Generation
		return true
	}

	return false
}

// watchACMECertificates watches filesystem for ACME certificate changes.
// v1 Conformance (INV-6.2): Event-driven ACME rotation via fsnotify.
func (r *CertReconciler) watchACMECertificates() {
	defer r.wg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.logger.Error("failed to create filesystem watcher", "err", err)
		return
	}
	defer watcher.Close()

	// Watch the directory containing ACME certificates
	certDir := filepath.Dir(r.acmeCertPath)
	if err := watcher.Add(certDir); err != nil {
		r.logger.Error("failed to watch ACME cert directory", "dir", certDir, "err", err)
		return
	}

	r.logger.Info("starting ACME certificate filesystem watcher",
		"cert", r.acmeCertPath,
		"key", r.acmeKeyPath)

	// Initialize current hashes
	r.initializeACMEHashes()

	// Debounce timer to avoid multiple rapid events
	var debounceTimer *time.Timer
	debounceDuration := 1 * time.Second

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("ACME certificate watcher stopped")
			return

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only care about Write and Create events for our cert/key files
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			if event.Name != r.acmeCertPath && event.Name != r.acmeKeyPath {
				continue
			}

			// Debounce: wait for file operations to complete
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceDuration, func() {
				if r.checkACMEHashChange() {
					// Signal ACME change (non-blocking)
					select {
					case r.acmeChanged <- struct{}{}:
						r.logger.Info("ACME certificate change detected", "file", event.Name)
					default:
						// Already pending, skip
					}
				}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			r.logger.Error("filesystem watch error", "err", err)
		}
	}
}

// initializeACMEHashes reads and stores initial ACME certificate hashes.
func (r *CertReconciler) initializeACMEHashes() {
	certHash := computeFileHash(r.acmeCertPath)
	keyHash := computeFileHash(r.acmeKeyPath)

	if certHash == "" || keyHash == "" {
		r.logger.Warn("ACME certificates not readable during initialization",
			"cert", r.acmeCertPath,
			"key", r.acmeKeyPath)
		return
	}

	r.mu.Lock()
	r.lastACMECertHash = certHash
	r.lastACMEKeyHash = keyHash
	r.mu.Unlock()

	r.logger.Info("ACME certificate hashes initialized")
}

// checkACMEHashChange checks if ACME certificate hashes have changed.
// Returns true if hashes changed and updates stored hashes.
func (r *CertReconciler) checkACMEHashChange() bool {
	certHash := computeFileHash(r.acmeCertPath)
	keyHash := computeFileHash(r.acmeKeyPath)

	// If files cannot be read, don't consider it a change
	if certHash == "" || keyHash == "" {
		return false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for changes
	if certHash != r.lastACMECertHash || keyHash != r.lastACMEKeyHash {
		r.lastACMECertHash = certHash
		r.lastACMEKeyHash = keyHash
		return true
	}

	return false
}
