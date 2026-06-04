// @awareness namespace=globular.platform
// @awareness component=platform_xds.cert_reconciler
// @awareness file_role=event_driven_cert_change_signaller_via_filesystem_watch
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=high
package watchers

import (
	"context"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// CertReconciler watches for certificate changes via the local
// filesystem and signals xDS to rebuild its snapshot.
//
// History: prior to v1.2.179 this reconciler watched the etcd key
// /globular/pki/bundles/{domain} directly via clientv3. That prefix
// is owned by node_agent's internal/certs package
// (services/golang/node_agent/node_agent_server/internal/certs/etcd_kv.go::PutBundle),
// so xDS reading raw etcd violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
//
// The migration switches to filesystem-based change detection: node-
// agent already writes the rotated cert / key / CA files to disk at
// /var/lib/globular/pki/issued/services/{service.crt, service.key}
// and /var/lib/globular/pki/ca.crt as part of its normal cert-issue
// flow. xDS watches those files via fsnotify and hashes their
// contents to suppress noise from atomic-rename `Create` events.
//
// Latency: comparable to the previous etcd Watch (sub-second). The
// fsnotify path has the additional bonus that an air-gapped install
// (no etcd routing for cert events) is unaffected.
type CertReconciler struct {
	logger *slog.Logger

	// Event channels
	certChanged chan struct{} // Signal channel for internal certificate changes
	acmeChanged chan struct{} // Signal channel for ACME certificate changes

	// State tracking
	mu                   sync.RWMutex
	lastInternalCertHash string
	lastInternalKeyHash  string
	lastInternalCAHash   string
	lastACMECertHash     string
	lastACMEKeyHash      string

	// Configuration — local filesystem paths only. No etcd, no
	// foreign-prefix reads.
	internalCertPath string // e.g. /var/lib/globular/pki/issued/services/service.crt
	internalKeyPath  string // e.g. /var/lib/globular/pki/issued/services/service.key
	internalCAPath   string // e.g. /var/lib/globular/pki/ca.crt
	acmeCertPath     string
	acmeKeyPath      string

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewCertReconciler creates a certificate reconciler with event-driven
// filesystem watching. All path arguments are optional: an empty
// internalCertPath/internalKeyPath/internalCAPath leaves the internal
// watcher idle; an empty acmeCertPath/acmeKeyPath leaves the ACME
// watcher idle.
func NewCertReconciler(logger *slog.Logger, internalCertPath, internalKeyPath, internalCAPath, acmeCertPath, acmeKeyPath string) *CertReconciler {
	if logger == nil {
		logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CertReconciler{
		logger:           logger,
		certChanged:      make(chan struct{}, 1),
		acmeChanged:      make(chan struct{}, 1),
		internalCertPath: internalCertPath,
		internalKeyPath:  internalKeyPath,
		internalCAPath:   internalCAPath,
		acmeCertPath:     acmeCertPath,
		acmeKeyPath:      acmeKeyPath,
		ctx:              ctx,
		cancel:           cancel,
	}
}

// Start begins watching for certificate changes.
func (r *CertReconciler) Start() error {
	// Start filesystem watcher for internal (cluster-CA-issued)
	// certificates. Requires at least one of the three internal
	// paths to be configured.
	if r.internalCertPath != "" || r.internalKeyPath != "" || r.internalCAPath != "" {
		r.wg.Add(1)
		go r.watchInternalCertificates()
	}

	// Start filesystem watcher for ACME certificates.
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

// watchInternalCertificates watches the local cert, key, and CA files
// for changes. Uses fsnotify on the containing directories (necessary
// for atomic-rename detection — writers typically write to a tmp file
// then `mv` over the target, which fsnotify reports as a Create on
// the target, not a Write).
//
// Hash comparison prevents spurious signals when the file is touched
// without content change. Debounce timer collapses rapid write/create
// pairs into a single signal.
func (r *CertReconciler) watchInternalCertificates() {
	defer r.wg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.logger.Error("failed to create internal-cert filesystem watcher", "err", err)
		return
	}
	defer watcher.Close()

	// Build the unique set of directories we need to watch — fsnotify
	// monitors directory entries, so any change to a file inside the
	// watched dir surfaces here.
	watched := make(map[string]bool)
	for _, p := range []string{r.internalCertPath, r.internalKeyPath, r.internalCAPath} {
		if p == "" {
			continue
		}
		dir := filepath.Dir(p)
		if watched[dir] {
			continue
		}
		if err := watcher.Add(dir); err != nil {
			r.logger.Warn("failed to watch internal-cert directory",
				"dir", dir, "err", err)
			continue
		}
		watched[dir] = true
	}
	if len(watched) == 0 {
		r.logger.Warn("internal-cert watcher started without any reachable directory — no signals will be emitted")
		return
	}

	r.logger.Info("starting internal-cert filesystem watcher",
		"cert", r.internalCertPath,
		"key", r.internalKeyPath,
		"ca", r.internalCAPath)

	r.initializeInternalHashes()

	const debounceDuration = 1 * time.Second
	var debounceTimer *time.Timer

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("internal-cert watcher stopped")
			return

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}
			// Filter to events on the paths we care about.
			if event.Name != r.internalCertPath &&
				event.Name != r.internalKeyPath &&
				event.Name != r.internalCAPath {
				continue
			}
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			triggerName := event.Name
			debounceTimer = time.AfterFunc(debounceDuration, func() {
				if r.checkInternalHashChange() {
					select {
					case r.certChanged <- struct{}{}:
						r.logger.Info("internal certificate change detected", "file", triggerName)
					default:
						// Already pending, skip.
					}
				}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			r.logger.Error("internal-cert filesystem watch error", "err", err)
		}
	}
}

// initializeInternalHashes computes and stores the initial hashes of
// the three internal cert files. Missing files leave the
// corresponding hash empty — checkInternalHashChange interprets an
// empty-to-nonempty transition as a change (first cert delivery).
func (r *CertReconciler) initializeInternalHashes() {
	certHash := computeFileHash(r.internalCertPath)
	keyHash := computeFileHash(r.internalKeyPath)
	caHash := computeFileHash(r.internalCAPath)

	r.mu.Lock()
	r.lastInternalCertHash = certHash
	r.lastInternalKeyHash = keyHash
	r.lastInternalCAHash = caHash
	r.mu.Unlock()

	r.logger.Info("internal-cert hashes initialized",
		"cert_present", certHash != "",
		"key_present", keyHash != "",
		"ca_present", caHash != "")
}

// checkInternalHashChange recomputes hashes and reports whether any
// changed. Updates stored hashes on change.
func (r *CertReconciler) checkInternalHashChange() bool {
	certHash := computeFileHash(r.internalCertPath)
	keyHash := computeFileHash(r.internalKeyPath)
	caHash := computeFileHash(r.internalCAPath)

	r.mu.Lock()
	defer r.mu.Unlock()

	changed := certHash != r.lastInternalCertHash ||
		keyHash != r.lastInternalKeyHash ||
		caHash != r.lastInternalCAHash

	if changed {
		r.lastInternalCertHash = certHash
		r.lastInternalKeyHash = keyHash
		r.lastInternalCAHash = caHash
	}
	return changed
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
