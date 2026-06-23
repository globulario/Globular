package globule

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// TestPrintBannerNeverWritesToStdout guards the describe contract:
//
// The node agent runs each service binary with --describe and parses stdout as
// JSON for port preflight. Any decorative/log output on stdout corrupts that
// JSON ("invalid character 'G' looking for beginning of value") and silently
// disables port preflight. PrintBanner is decorative human output, so it MUST
// go to stderr only.
//
// Regression for: gateway/xds --describe emitting the "Globular v..." banner on
// stdout. See invariant observability.describe_stdout_is_machine_only.
func TestPrintBannerNeverWritesToStdout(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	PrintBanner("9.9.9", 42)

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = origStdout

	captured, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}
	if len(bytes.TrimSpace(captured)) != 0 {
		t.Fatalf("PrintBanner wrote to stdout (must be stderr-only); got %q", captured)
	}
}
