package bootstrap

import (
	"context"
	"io"
	"os"

	"github.com/globulario/services/golang/process"
)

// StartDNSService starts the DNS service and returns the PID.
func StartDNSService(ctx context.Context, cfg map[string]any, port int, outW, errW io.Writer) (int, error) {
	if outW == nil {
		outW = os.Stdout
	}
	if errW == nil {
		errW = os.Stderr
	}
	return process.StartServiceProcessWithWriters(cfg, port, outW, errW)
}

// StartDNSProxy starts the proxy for the DNS service.
func StartDNSProxy(ctx context.Context, cfg map[string]any, caPath, certPath string) (int, error) {
	return process.StartServiceProxyProcess(cfg, caPath, certPath)
}
