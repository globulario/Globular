package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// binarySpec maps a URL name to its binary filename on disk.
type binarySpec struct {
	BinaryName string // filename in binDir
}

var allowedBinaries = map[string]binarySpec{
	"node_agent_server": {BinaryName: "node_agent_server"},
	"globularcli":       {BinaryName: "globularcli"},
	"etcd":              {BinaryName: "etcd"},
	"etcdctl":           {BinaryName: "etcdctl"},
}

// NewJoinBinHandler serves binaries for joining nodes directly from disk.
// binDir is the directory containing the compiled binaries (e.g. /usr/lib/globular/bin).
func NewJoinBinHandler(binDir string) http.Handler {
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

		if _, ok := allowedBinaries[name]; !ok {
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

		path := filepath.Join(binDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, "binary not found", http.StatusNotFound)
			return
		}

		cacheMu.Lock()
		cache[name] = cachedBin{data: data, at: time.Now()}
		cacheMu.Unlock()

		serveBinaryBytes(w, name, data)
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

// NewJoinPackagesHandler serves package .tgz files from pkgDir.
// GET /join/packages/         → JSON list of available package filenames
// GET /join/packages/<file>   → streams the .tgz file
func NewJoinPackagesHandler(pkgDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tail := strings.TrimPrefix(r.URL.Path, "/join/packages/")
		tail = strings.Trim(tail, "/")

		if tail == "" {
			// List all .tgz files.
			entries, err := os.ReadDir(pkgDir)
			if err != nil {
				if os.IsNotExist(err) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("[]"))
					return
				}
				http.Error(w, "list packages: "+err.Error(), http.StatusInternalServerError)
				return
			}
			var names []string
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".tgz") {
					names = append(names, e.Name())
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(names)
			return
		}

		// Serve a specific file — sanitize to prevent directory traversal.
		name := filepath.Base(tail)
		if !strings.HasSuffix(name, ".tgz") {
			http.Error(w, "only .tgz files are served here", http.StatusForbidden)
			return
		}
		path := filepath.Join(pkgDir, name)
		if _, err := os.Stat(path); err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		http.ServeFile(w, r, path)
	})
}
