package bootstrap

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// StartMinio launches the MinIO binary with the provided configuration.
func StartMinio(ctx context.Context, bin, listenAddr, certsDir, dataDir string, env []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	cmd := exec.CommandContext(ctx, bin,
		"server",
		"--certs-dir", certsDir,
		"--address", listenAddr,
		dataDir,
	)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("startMinio: %w", err)
	}

	go func() {
		_ = cmd.Wait()
	}()

	// Give minio a moment to settle before returning
	time.Sleep(2 * time.Second)
	return cmd, nil
}
