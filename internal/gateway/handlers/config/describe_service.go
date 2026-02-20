package config

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/globulario/services/golang/config"
)

// ServiceDescriber fetches the --describe metadata for a service binary.
type ServiceDescriber interface {
	DescribeService(name string, timeout time.Duration) (config.ServiceDesc, string, error)
}

// NewDescribeService returns a handler for GET /api/describe-service?name=<service>&timeoutMs=<n>.
func NewDescribeService(sd ServiceDescriber) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceName := strings.TrimSpace(r.URL.Query().Get("name"))
		if serviceName == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}

		timeout := parseTimeout(r.URL.Query().Get("timeoutMs"))
		desc, path, err := sd.DescribeService(serviceName, timeout)
		if err != nil {
			http.Error(w, "describe failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		resp := struct {
			config.ServiceDesc
			BinaryPath string `json:"binaryPath"`
		}{
			ServiceDesc: desc,
			BinaryPath:  path,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

func parseTimeout(val string) time.Duration {
	if ms, err := strconv.Atoi(strings.TrimSpace(val)); err == nil && ms > 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return 0
}
