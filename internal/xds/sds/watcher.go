package sds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Watcher monitors certificate generation changes in etcd and triggers SDS updates.
// When a certificate is rotated (generation increments), it reloads secrets from disk
// and pushes updated secrets to Envoy via the SDS server.
type Watcher struct {
	etcdClient *clientv3.Client
	sdsServer  *Server
	domain     string

	// Certificate file paths (read from disk when generation changes)
	internalCertPath string
	internalKeyPath  string
	caPath           string
	publicCertPath   string
	publicKeyPath    string

	// State tracking
	lastGeneration uint64
	pollInterval   time.Duration

	// Interfaces for testability
	certKV certKV
}

// certKV abstracts the etcd certificate storage interface
type certKV interface {
	GetBundleGeneration(ctx context.Context, domain string) (uint64, error)
}

// etcdCertKV implements certKV using etcd directly
type etcdCertKV struct {
	client *clientv3.Client
}

func (kv *etcdCertKV) GetBundleGeneration(ctx context.Context, domain string) (uint64, error) {
	// Use same key pattern as services/golang/nodeagent
	key := fmt.Sprintf("/globular/pki/bundles/%s", domain)

	resp, err := kv.client.Get(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("etcd get: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return 0, fmt.Errorf("bundle not found for domain %s", domain)
	}

	// Parse JSON to extract generation
	// Format: {"generation": 123, "updated_ms": ..., "key": "...", "fullchain": "...", "ca": "..."}
	var payload struct {
		Generation uint64 `json:"generation"`
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &payload); err != nil {
		return 0, fmt.Errorf("parse bundle: %w", err)
	}

	return payload.Generation, nil
}

// WatcherConfig holds configuration for the certificate rotation watcher.
type WatcherConfig struct {
	EtcdClient *clientv3.Client
	SDSServer  *Server
	Domain     string

	// Certificate file paths (canonical PKI locations)
	InternalCertPath string
	InternalKeyPath  string
	CAPath           string

	// Optional public cert paths (ACME)
	PublicCertPath string
	PublicKeyPath  string

	// Poll interval (default 10s)
	PollInterval time.Duration
}

// NewWatcher creates a new certificate rotation watcher.
func NewWatcher(cfg WatcherConfig) (*Watcher, error) {
	if cfg.EtcdClient == nil {
		return nil, fmt.Errorf("etcd client is required")
	}
	if cfg.SDSServer == nil {
		return nil, fmt.Errorf("SDS server is required")
	}
	if cfg.Domain == "" {
		return nil, fmt.Errorf("domain is required")
	}
	if cfg.InternalCertPath == "" || cfg.InternalKeyPath == "" || cfg.CAPath == "" {
		return nil, fmt.Errorf("internal cert paths are required")
	}

	pollInterval := cfg.PollInterval
	if pollInterval == 0 {
		pollInterval = 10 * time.Second
	}

	return &Watcher{
		etcdClient:       cfg.EtcdClient,
		sdsServer:        cfg.SDSServer,
		domain:           cfg.Domain,
		internalCertPath: cfg.InternalCertPath,
		internalKeyPath:  cfg.InternalKeyPath,
		caPath:           cfg.CAPath,
		publicCertPath:   cfg.PublicCertPath,
		publicKeyPath:    cfg.PublicKeyPath,
		pollInterval:     pollInterval,
		certKV:           &etcdCertKV{client: cfg.EtcdClient},
	}, nil
}

// Run starts the watcher loop. It polls etcd at PollInterval and checks for generation changes.
// When a change is detected, it reloads secrets from disk and updates the SDS server.
// Returns when ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	log.Printf("sds watcher: starting for domain %s (poll interval: %s)", w.domain, w.pollInterval)

	// Load initial secrets
	if err := w.checkAndUpdateSecrets(ctx); err != nil {
		log.Printf("sds watcher: initial load: %v", err)
	}

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("sds watcher: stopped")
			return ctx.Err()

		case <-ticker.C:
			if err := w.checkAndUpdateSecrets(ctx); err != nil {
				log.Printf("sds watcher: check failed: %v", err)
			}
		}
	}
}

// checkAndUpdateSecrets checks the certificate generation in etcd.
// If changed, reloads secrets from disk and updates the SDS server.
func (w *Watcher) checkAndUpdateSecrets(ctx context.Context) error {
	// Get current generation from etcd
	gen, err := w.certKV.GetBundleGeneration(ctx, w.domain)
	if err != nil {
		return fmt.Errorf("get generation: %w", err)
	}

	// Check if generation changed
	if gen == w.lastGeneration && w.lastGeneration != 0 {
		// No change
		return nil
	}

	log.Printf("sds watcher: generation changed: %d -> %d", w.lastGeneration, gen)

	// Reload secrets from disk
	secrets, err := w.buildSecrets()
	if err != nil {
		return fmt.Errorf("build secrets: %w", err)
	}

	// Update SDS server (this pushes to Envoy via snapshot cache)
	if err := w.sdsServer.UpdateSecrets(secrets); err != nil {
		return fmt.Errorf("update secrets: %w", err)
	}

	w.lastGeneration = gen
	log.Printf("sds watcher: secrets updated for generation %d", gen)

	return nil
}

// buildSecrets reads certificate files from disk and builds Secret resources.
func (w *Watcher) buildSecrets() (map[string]*tls_v3.Secret, error) {
	internalPaths := CertPaths{
		CertFile: w.internalCertPath,
		KeyFile:  w.internalKeyPath,
	}

	publicPaths := CertPaths{
		CertFile: w.publicCertPath,
		KeyFile:  w.publicKeyPath,
	}

	return BuildAllSecrets(internalPaths, publicPaths, w.caPath)
}

// GetLastGeneration returns the last observed certificate generation.
// Useful for testing and debugging.
func (w *Watcher) GetLastGeneration() uint64 {
	return w.lastGeneration
}
