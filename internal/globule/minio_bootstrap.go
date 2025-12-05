package globule

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/globulario/services/golang/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioRuntime captures the runtime details Globule needs to configure FileService.
type MinioRuntime struct {
	Endpoint   string // host:port for SDK usage
	PublicURL  string // optional https:// URL for logging
	AccessKey  string
	SecretKey  string
	Bucket     string
	Prefix     string
	UseSSL     bool
	CertsDir   string
	DataDir    string
	ListenAddr string
}

// prepareMinioCerts creates a certs dir for MinIO and links the existing Globule certs.
func (g *Globule) prepareMinioCerts() (string, error) {
	pub := config.GetLocalCertificate()
	key := config.GetLocalServerKeyPath()

	if _, err := os.Stat(pub); err != nil {
		return "", fmt.Errorf("prepareMinioCerts: server certificate not found at %s: %w", pub, err)
	}
	if _, err := os.Stat(key); err != nil {
		return "", fmt.Errorf("prepareMinioCerts: server key not found at %s: %w", key, err)
	}

	baseDir := filepath.Dir(pub)
	certsDir := filepath.Join(baseDir, "minio-certs")
	if err := os.MkdirAll(certsDir, 0o700); err != nil {
		return "", fmt.Errorf("prepareMinioCerts: mkdir %s: %w", certsDir, err)
	}

	dstPub := filepath.Join(certsDir, "public.crt")
	dstKey := filepath.Join(certsDir, "private.key")

	link := func(target, linkName string) error {
		_ = os.Remove(linkName)
		if err := os.Symlink(target, linkName); err != nil {
			if !os.IsPermission(err) {
				return err
			}
			data, readErr := os.ReadFile(target)
			if readErr != nil {
				return readErr
			}
			return os.WriteFile(linkName, data, 0o600)
		}
		return nil
	}

	if err := link(pub, dstPub); err != nil {
		return "", fmt.Errorf("prepareMinioCerts: link public: %w", err)
	}
	if err := link(key, dstKey); err != nil {
		return "", fmt.Errorf("prepareMinioCerts: link private: %w", err)
	}

	return certsDir, nil
}

// startMinioIfNeeded boots MinIO unless MINIO_DISABLED=1, wiring TLS + credentials.
func (g *Globule) startMinioIfNeeded(ctx context.Context, log *slog.Logger) (*MinioRuntime, error) {
	if v := os.Getenv("MINIO_DISABLED"); v == "1" {
		return nil, nil
	}

	accessKey := firstNonEmpty(strings.TrimSpace(os.Getenv("MINIO_ROOT_USER")), "globular")
	secretKey := firstNonEmpty(strings.TrimSpace(os.Getenv("MINIO_ROOT_PASSWORD")), "globular-secret")
	bucket := firstNonEmpty(strings.TrimSpace(os.Getenv("MINIO_BUCKET")), "globular")
	prefix := strings.TrimSpace(os.Getenv("MINIO_PREFIX"))

	dataDir := firstNonEmpty(strings.TrimSpace(os.Getenv("MINIO_DATA_DIR")), "/mnt/globular-minio")
	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return nil, fmt.Errorf("startMinio: mkdir dataDir %s: %w", dataDir, err)
	}

	host := config.HostOnly(g.localIPAddress)
	if strings.TrimSpace(host) == "" {
		host = "0.0.0.0"
	}
	port := firstNonEmpty(strings.TrimSpace(os.Getenv("MINIO_PORT")), "9000")
	listenAddr := host + ":" + port

	endpointHost := host
	if d, err := config.GetDomain(); err == nil && strings.TrimSpace(d) != "" {
		endpointHost = strings.TrimSpace(d)
	}
	clientEndpoint := endpointHost + ":" + port
	publicURL := "https://" + clientEndpoint

	certsDir, err := g.prepareMinioCerts()
	if err != nil {
		return nil, err
	}

	minioBin := firstNonEmpty(strings.TrimSpace(os.Getenv("MINIO_BIN")), "minio")
	cmd := exec.CommandContext(ctx, minioBin,
		"server",
		"--certs-dir", certsDir,
		"--address", listenAddr,
		dataDir,
	)
	cmd.Env = append(os.Environ(),
		"MINIO_ROOT_USER="+accessKey,
		"MINIO_ROOT_PASSWORD="+secretKey,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Info("starting MinIO",
		"bin", minioBin,
		"addr", listenAddr,
		"publicURL", publicURL,
		"dataDir", dataDir,
		"certsDir", certsDir,
	)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("startMinio: %w", err)
	}

	go func() {
		_ = cmd.Wait()
		log.Warn("MinIO process exited")
	}()

	time.Sleep(2 * time.Second)

	rt := &MinioRuntime{
		Endpoint:   clientEndpoint,
		PublicURL:  publicURL,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Bucket:     bucket,
		Prefix:     prefix,
		UseSSL:     true,
		CertsDir:   certsDir,
		DataDir:    dataDir,
		ListenAddr: listenAddr,
	}

	if err := g.ensureMinioBucket(ctx, log, rt); err != nil {
		return nil, err
	}

	return rt, nil
}

func (g *Globule) ensureMinioBucket(ctx context.Context, log *slog.Logger, cfg *MinioRuntime) error {
	bucket := strings.TrimSpace(cfg.Bucket)
	if bucket == "" {
		return fmt.Errorf("minio bucket name is empty")
	}

	minioClient, err := g.newMinioClient(log, cfg)
	if err != nil {
		return fmt.Errorf("failed to init MinIO client: %w", err)
	}

	const maxAttempts = 15
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("ensure MinIO bucket %s: %w", bucket, ctx.Err())
		default:
		}

		exists, err := minioClient.BucketExists(ctx, bucket)
		if err == nil {
			if exists {
				log.Info("MinIO bucket ready", "bucket", bucket)
				return nil
			}
			if err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err == nil {
				log.Info("created MinIO bucket", "bucket", bucket)
				return nil
			}
		}

		lastErr = err
		log.Warn("waiting for MinIO bucket", "bucket", bucket, "attempt", attempt, "err", err)

		select {
		case <-ctx.Done():
			return fmt.Errorf("ensure MinIO bucket %s: %w", bucket, ctx.Err())
		case <-time.After(2 * time.Second):
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("bucket %s not ready", bucket)
	}
	return fmt.Errorf("ensure MinIO bucket %s: %w", bucket, lastErr)
}

func (g *Globule) newMinioClient(log *slog.Logger, cfg *MinioRuntime) (*minio.Client, error) {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	}

	if cfg.UseSSL {
		transport := cloneHTTPTransport()
		transport.TLSClientConfig = g.minioTLSConfig(log)
		opts.Transport = transport
	}

	return minio.New(cfg.Endpoint, opts)
}

func cloneHTTPTransport() *http.Transport {
	if base, ok := http.DefaultTransport.(*http.Transport); ok {
		return base.Clone()
	}
	return &http.Transport{}
}

func (g *Globule) minioTLSConfig(log *slog.Logger) *tls.Config {
	cfg := &tls.Config{}

	caPath := strings.TrimSpace(config.GetLocalCertificateAuthorityBundle())
	if caPath == "" {
		cfg.InsecureSkipVerify = true
		log.Warn("no CA certificate configured for MinIO; skipping TLS verification")
		return cfg
	}

	data, err := os.ReadFile(caPath)
	if err != nil {
		cfg.InsecureSkipVerify = true
		log.Warn("failed to read CA certificate for MinIO; skipping TLS verification", "path", caPath, "err", err)
		return cfg
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		cfg.InsecureSkipVerify = true
		log.Warn("failed to parse CA certificate for MinIO; skipping TLS verification", "path", caPath)
		return cfg
	}

	cfg.RootCAs = pool
	return cfg
}
