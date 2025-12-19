package lifecycle

import (
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/globulario/Globular/internal/legacy"
	"github.com/globulario/services/golang/process"
)

// Lifecycle provides helper methods that execute OS-level service commands.
type Lifecycle struct {
	deps legacy.Deps
}

// New creates a Lifecycle handler using the supplied dependencies.
func New(deps legacy.Deps) *Lifecycle {
	if deps.Timeout <= 0 {
		deps.Timeout = 1500 * time.Millisecond
	}
	if deps.Stdout == nil {
		deps.Stdout = os.Stdout
	}
	if deps.Stderr == nil {
		deps.Stderr = os.Stderr
	}
	return &Lifecycle{deps: deps}
}

// StartService starts the given service using the configured process helper.
func (l *Lifecycle) StartService(ctx context.Context, cfg map[string]any, port int, outW, errW io.Writer) (int, error) {
	if outW == nil {
		outW = l.deps.Stdout
	}
	if errW == nil {
		errW = l.deps.Stderr
	}
	return process.StartServiceProcessWithWriters(cfg, port, outW, errW)
}

// StartProxy creates a service proxy for the provided configuration.
func (l *Lifecycle) StartProxy(ctx context.Context, cfg map[string]any, caFile, certFile string) (int, error) {
	return process.StartServiceProxyProcess(cfg, caFile, certFile)
}

// RunHealthCheck executes the provided binary with the health args and env.
func (l *Lifecycle) RunHealthCheck(ctx context.Context, bin string, args, env []string) error {
	ctx, cancel := context.WithTimeout(ctx, l.deps.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = l.deps.Stdout
	cmd.Stderr = l.deps.Stderr
	return cmd.Run()
}

// KillService stops a running service process.
func (l *Lifecycle) KillService(cfg map[string]any) error {
	return process.KillServiceProcess(cfg)
}

// KillProxy stops a service proxy process.
func (l *Lifecycle) KillProxy(cfg map[string]any) error {
	return process.KillServiceProxyProcess(cfg)
}
