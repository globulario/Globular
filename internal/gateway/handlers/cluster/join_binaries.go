package cluster

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	repopb "github.com/globulario/services/golang/repository/repositorypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// binarySpec maps a URL name to its repository package and binary filename.
type binarySpec struct {
	Package    string // repository package name (e.g. "node-agent")
	BinaryName string // filename inside bin/ in the tarball
}

var allowedBinaries = map[string]binarySpec{
	"node_agent_server": {Package: "node-agent", BinaryName: "node_agent_server"},
	"globularcli":       {Package: "globular-cli", BinaryName: "globularcli"},
	"etcd":              {Package: "etcd", BinaryName: "etcd"},
	"etcdctl":           {Package: "etcdctl", BinaryName: "etcdctl"},
}

// NewJoinBinHandler serves binaries for joining nodes by downloading them
// from the repository service. This guarantees the served binary has the
// same digest as the published artifact — no self-upgrade mismatch.
//
// Falls back to disk (binDir) if the repository is unavailable.
// NewJoinBinHandler serves join binaries. repoAddr is the repository service
// endpoint resolved from etcd (e.g. "10.0.0.63:10003") — never localhost.
func NewJoinBinHandler(binDir, repoAddr string) http.Handler {
	var (
		cache   = make(map[string]cachedBin)
		cacheMu sync.RWMutex
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/join/bin/")
		name = filepath.Base(name)

		spec, ok := allowedBinaries[name]
		if !ok {
			http.NotFound(w, r)
			return
		}

		// Check cache (5 min TTL).
		cacheMu.RLock()
		cached, hit := cache[name]
		cacheMu.RUnlock()
		if hit && time.Since(cached.at) < 5*time.Minute && len(cached.data) > 0 {
			serveBinaryBytes(w, name, cached.data)
			return
		}

		// Download from repository.
		data, err := fetchBinaryFromRepo(spec, repoAddr)
		if err == nil && len(data) > 0 {
			cacheMu.Lock()
			cache[name] = cachedBin{data: data, at: time.Now()}
			cacheMu.Unlock()
			serveBinaryBytes(w, name, data)
			return
		}
		log.Printf("join/bin/%s: repo fetch failed: %v — falling back to disk", name, err)

		// Fallback to disk.
		path := filepath.Join(binDir, name)
		if _, err := os.Stat(path); err != nil {
			http.Error(w, "binary not found in repository or on disk", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		http.ServeFile(w, r, path)
	})
}

type cachedBin struct {
	data []byte
	at   time.Time
}

func serveBinaryBytes(w http.ResponseWriter, name string, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write(data)
}

// fetchBinaryFromRepo downloads the package .tgz from the repository via
// direct gRPC (bypassing Envoy mesh), extracts the requested binary, and
// returns its contents. This is the single source of truth for join delivery.
func fetchBinaryFromRepo(spec binarySpec, repoAddr string) ([]byte, error) {
	conn, err := directRepoConn(repoAddr)
	if err != nil {
		return nil, fmt.Errorf("repo connection: %w", err)
	}
	defer conn.Close()

	client := repopb.NewPackageRepositoryClient(conn)

	platform := runtime.GOOS + "_" + runtime.GOARCH

	// Resolve the exact package ref from the active platform release BOM rather
	// than blindly fetching "latest". If the cluster has synced v1.0.85 but the
	// active release is v1.0.84, joining nodes must get v1.0.84 binaries.
	//
	// Resolution chain:
	//   1. Read /var/lib/globular/release-index.json (written at Day-0)
	//   2. Look up version + build_number for this spec.Package
	//   3. Fallback to latest published if no BOM exists (legacy mode)
	bom := resolveFromBOM(spec.Package, platform)

	publisher := bom.Publisher
	if publisher == "" {
		publisher = "core@globular.io"
	}
	ref := &repopb.ArtifactRef{
		Name:        spec.Package,
		Version:     bom.Version,
		Platform:    platform,
		PublisherId: publisher,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := client.DownloadArtifact(ctx, &repopb.DownloadArtifactRequest{
		Ref:         ref,
		BuildNumber: bom.BuildNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", spec.Package, err)
	}

	// Collect all chunks into a buffer.
	var buf bytes.Buffer
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("receive chunk for %s: %w", spec.Package, err)
		}
		buf.Write(resp.GetData())
	}

	// Extract the binary from the tarball.
	gz, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tar: %w", err)
		}
		if hdr.FileInfo().IsDir() {
			continue
		}
		if filepath.Base(hdr.Name) != spec.BinaryName {
			continue
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", spec.BinaryName, err)
		}
		log.Printf("join/bin: extracted %s (%d bytes) from repo artifact %s",
			spec.BinaryName, len(data), spec.Package)
		return data, nil
	}

	return nil, fmt.Errorf("binary %s not found in artifact %s", spec.BinaryName, spec.Package)
}

// bomPackageRef holds the resolved package identity from the active BOM.
type bomPackageRef struct {
	Version     string
	BuildNumber int64
	Publisher   string
}

// resolveFromBOM reads the active release-index.json (written at Day-0) and
// returns the exact package ref for the requested binary. Includes version
// AND build_number for precise artifact resolution.
// If no BOM exists (legacy), returns zero-value ref (version="" → latest).
func resolveFromBOM(pkgName, platform string) bomPackageRef {
	const releaseIndexPath = "/var/lib/globular/release-index.json"

	data, err := os.ReadFile(releaseIndexPath)
	if err != nil {
		// No BOM — legacy mode, fall back to latest published.
		return bomPackageRef{}
	}

	var idx struct {
		Packages []struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			BuildNumber int64  `json:"build_number"`
			Publisher   string `json:"publisher"`
			Platform    string `json:"platform"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(data, &idx); err != nil {
		log.Printf("join/bin: failed to parse release-index.json: %v — using latest", err)
		return bomPackageRef{}
	}

	for _, p := range idx.Packages {
		if p.Name == pkgName && (platform == "" || p.Platform == platform) {
			pub := p.Publisher
			if pub == "" {
				pub = "core@globular.io"
			}
			log.Printf("join/bin: resolved %s v%s build=%d from active release BOM",
				pkgName, p.Version, p.BuildNumber)
			return bomPackageRef{
				Version:     p.Version,
				BuildNumber: p.BuildNumber,
				Publisher:   pub,
			}
		}
	}

	log.Printf("join/bin: package %s not found in release-index.json — using latest", pkgName)
	return bomPackageRef{}
}

// directRepoConn establishes a direct mTLS gRPC connection to the repository
// service at the given address (resolved from etcd by the caller — never
// localhost or a hardcoded port).
func directRepoConn(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, fmt.Errorf("repository address not provided (etcd resolution failed)")
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid repository address %q: %w", addr, err)
	}

	certFile := "/var/lib/globular/pki/issued/services/service.crt"
	keyFile := "/var/lib/globular/pki/issued/services/service.key"
	caFile := "/var/lib/globular/pki/ca.crt"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("read CA: %w", err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caPEM)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   host, // matches the cert SAN (node's actual IP/hostname)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds))
}
