package sds

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// mockCertKV implements certKV for testing
type mockCertKV struct {
	generation uint64
	err        error
}

func (m *mockCertKV) GetBundleGeneration(ctx context.Context, domain string) (uint64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.generation, nil
}

func TestNewWatcher(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	// Create mock etcd client (nil is OK for test, we'll inject mock certKV)
	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{}, // Mock client
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
	}

	watcher, err := NewWatcher(cfg)
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}

	if watcher.domain != "globular.internal" {
		t.Errorf("expected domain globular.internal, got %s", watcher.domain)
	}

	if watcher.pollInterval != 10*time.Second {
		t.Errorf("expected default poll interval 10s, got %s", watcher.pollInterval)
	}
}

func TestNewWatcher_CustomPollInterval(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
		PollInterval:     5 * time.Second,
	}

	watcher, err := NewWatcher(cfg)
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}

	if watcher.pollInterval != 5*time.Second {
		t.Errorf("expected poll interval 5s, got %s", watcher.pollInterval)
	}
}

func TestNewWatcher_MissingEtcdClient(t *testing.T) {
	cfg := WatcherConfig{
		SDSServer: NewServer(),
		Domain:    "globular.internal",
	}

	_, err := NewWatcher(cfg)
	if err == nil {
		t.Error("expected error for missing etcd client")
	}
}

func TestNewWatcher_MissingSDSServer(t *testing.T) {
	cfg := WatcherConfig{
		EtcdClient: &clientv3.Client{},
		Domain:     "globular.internal",
	}

	_, err := NewWatcher(cfg)
	if err == nil {
		t.Error("expected error for missing SDS server")
	}
}

func TestNewWatcher_MissingDomain(t *testing.T) {
	cfg := WatcherConfig{
		EtcdClient: &clientv3.Client{},
		SDSServer:  NewServer(),
	}

	_, err := NewWatcher(cfg)
	if err == nil {
		t.Error("expected error for missing domain")
	}
}

func TestNewWatcher_MissingCertPaths(t *testing.T) {
	cfg := WatcherConfig{
		EtcdClient: &clientv3.Client{},
		SDSServer:  NewServer(),
		Domain:     "globular.internal",
	}

	_, err := NewWatcher(cfg)
	if err == nil {
		t.Error("expected error for missing cert paths")
	}
}

func TestCheckAndUpdateSecrets_InitialLoad(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
	}

	watcher, _ := NewWatcher(cfg)

	// Inject mock certKV
	mockKV := &mockCertKV{generation: 1}
	watcher.certKV = mockKV

	// Initial load (generation 0 -> 1)
	err := watcher.checkAndUpdateSecrets(context.Background())
	if err != nil {
		t.Fatalf("checkAndUpdateSecrets failed: %v", err)
	}

	if watcher.GetLastGeneration() != 1 {
		t.Errorf("expected generation 1, got %d", watcher.GetLastGeneration())
	}

	// Verify secrets loaded into server
	if len(srv.secrets) == 0 {
		t.Error("expected secrets to be loaded into server")
	}
}

func TestCheckAndUpdateSecrets_NoChange(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
	}

	watcher, _ := NewWatcher(cfg)
	watcher.lastGeneration = 5

	// Inject mock certKV with same generation
	mockKV := &mockCertKV{generation: 5}
	watcher.certKV = mockKV

	// Should not update (generation unchanged)
	err := watcher.checkAndUpdateSecrets(context.Background())
	if err != nil {
		t.Fatalf("checkAndUpdateSecrets failed: %v", err)
	}

	// Verify server secrets not updated
	if len(srv.secrets) != 0 {
		t.Error("expected no secrets update when generation unchanged")
	}
}

func TestCheckAndUpdateSecrets_GenerationIncrement(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
	}

	watcher, _ := NewWatcher(cfg)
	watcher.lastGeneration = 5

	// Inject mock certKV with incremented generation
	mockKV := &mockCertKV{generation: 6}
	watcher.certKV = mockKV

	initialVersion := srv.GetVersion()

	// Should update (generation changed)
	err := watcher.checkAndUpdateSecrets(context.Background())
	if err != nil {
		t.Fatalf("checkAndUpdateSecrets failed: %v", err)
	}

	if watcher.GetLastGeneration() != 6 {
		t.Errorf("expected generation 6, got %d", watcher.GetLastGeneration())
	}

	// Verify server secrets updated
	if len(srv.secrets) == 0 {
		t.Error("expected secrets to be loaded")
	}

	// Verify version changed
	if srv.GetVersion() == initialVersion {
		t.Error("expected version to change after update")
	}
}

func TestCheckAndUpdateSecrets_EtcdError(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
	}

	watcher, _ := NewWatcher(cfg)

	// Inject mock certKV with error
	mockKV := &mockCertKV{err: fmt.Errorf("etcd connection failed")}
	watcher.certKV = mockKV

	// Should return error
	err := watcher.checkAndUpdateSecrets(context.Background())
	if err == nil {
		t.Error("expected error from etcd")
	}
}

func TestCheckAndUpdateSecrets_MissingFiles(t *testing.T) {
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: "/nonexistent/cert.pem",
		InternalKeyPath:  "/nonexistent/key.pem",
		CAPath:           "/nonexistent/ca.pem",
	}

	watcher, _ := NewWatcher(cfg)

	// Inject mock certKV
	mockKV := &mockCertKV{generation: 1}
	watcher.certKV = mockKV

	// Should fail to build secrets
	err := watcher.checkAndUpdateSecrets(context.Background())
	if err == nil {
		t.Error("expected error for missing cert files")
	}
}

func TestBuildSecrets(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
	}

	watcher, _ := NewWatcher(cfg)

	secrets, err := watcher.buildSecrets()
	if err != nil {
		t.Fatalf("buildSecrets failed: %v", err)
	}

	// Should have internal-server-cert and internal-ca-bundle
	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}

	if _, ok := secrets[InternalServerCert]; !ok {
		t.Error("missing internal-server-cert")
	}

	if _, ok := secrets[InternalCABundle]; !ok {
		t.Error("missing internal-ca-bundle")
	}
}

func TestBuildSecrets_WithPublicCert(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	publicCertPath := tmpDir + "/public-cert.pem"
	publicKeyPath := tmpDir + "/public-key.pem"

	if err := os.WriteFile(publicCertPath, []byte(testCert), 0644); err != nil {
		t.Fatalf("write public cert: %v", err)
	}
	if err := os.WriteFile(publicKeyPath, []byte(testKey), 0600); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
		PublicCertPath:   publicCertPath,
		PublicKeyPath:    publicKeyPath,
	}

	watcher, _ := NewWatcher(cfg)

	secrets, err := watcher.buildSecrets()
	if err != nil {
		t.Fatalf("buildSecrets failed: %v", err)
	}

	// Should have 3 secrets now
	if len(secrets) != 3 {
		t.Errorf("expected 3 secrets, got %d", len(secrets))
	}

	if _, ok := secrets[PublicServerCert]; !ok {
		t.Error("missing public-server-cert")
	}
}

func TestWatcherRun_Cancellation(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
		PollInterval:     100 * time.Millisecond, // Fast polling for test
	}

	watcher, _ := NewWatcher(cfg)

	// Inject mock certKV
	mockKV := &mockCertKV{generation: 1}
	watcher.certKV = mockKV

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Run should return when context cancelled
	err := watcher.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("expected context error, got %v", err)
	}

	// Verify watcher ran at least once
	if watcher.GetLastGeneration() != 1 {
		t.Error("watcher should have run at least once")
	}
}

func TestWatcherRun_MultiplePolls(t *testing.T) {
	tmpDir := t.TempDir()

	certPath, keyPath, caPath := createTestCertFiles(t, tmpDir)
	srv := NewServer()

	cfg := WatcherConfig{
		EtcdClient:       &clientv3.Client{},
		SDSServer:        srv,
		Domain:           "globular.internal",
		InternalCertPath: certPath,
		InternalKeyPath:  keyPath,
		CAPath:           caPath,
		PollInterval:     50 * time.Millisecond,
	}

	watcher, _ := NewWatcher(cfg)

	// Inject mock certKV that increments generation over time
	mockKV := &mockCertKV{generation: 1}
	watcher.certKV = mockKV

	ctx, cancel := context.WithCancel(context.Background())

	// Run in background
	done := make(chan error)
	go func() {
		done <- watcher.Run(ctx)
	}()

	// Wait for initial poll
	time.Sleep(100 * time.Millisecond)

	// Change generation
	mockKV.generation = 2

	// Wait for next poll
	time.Sleep(100 * time.Millisecond)

	// Cancel and wait
	cancel()
	<-done

	// Verify generation updated
	if watcher.GetLastGeneration() != 2 {
		t.Errorf("expected generation 2, got %d", watcher.GetLastGeneration())
	}
}

func TestEtcdCertKV_GetBundleGeneration(t *testing.T) {
	// This test requires a real etcd instance, skip if not available
	t.Skip("Integration test - requires etcd")

	// Example of how to test with real etcd:
	// client, _ := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	// kv := &etcdCertKV{client: client}
	// gen, err := kv.GetBundleGeneration(context.Background(), "globular.internal")
	// ...
}
