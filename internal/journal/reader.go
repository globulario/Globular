// Package journal provides helpers to read systemd journal entries
// via journalctl. This package is intentionally outside internal/gateway
// so it can use os/exec without violating the gateway no-exec lint rule.
package journal

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const defaultTimeout = 3 * time.Second

// Result holds the output of a journal read.
type Result struct {
	Unit      string
	Lines     []string
	Truncated bool
	Error     string
}

// ReadUnit reads the last N lines from a systemd unit's journal.
// sinceSec limits the time window (e.g. 3600 = last hour).
// The context should carry any caller-imposed deadline; an additional
// internal timeout of 3s is applied to the journalctl invocation.
func ReadUnit(ctx context.Context, unit string, lines int, sinceSec int) Result {
	if lines < 1 {
		lines = 1
	}
	if lines > 500 {
		lines = 500
	}

	execCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	sinceArg := fmt.Sprintf("%d seconds ago", sinceSec)
	cmd := exec.CommandContext(execCtx, "journalctl",
		"-u", unit,
		"-n", strconv.Itoa(lines),
		"--no-pager",
		"-o", "short-iso",
		"--since", sinceArg,
	)
	out, err := cmd.CombinedOutput()

	r := Result{Unit: unit}
	if err != nil {
		r.Error = err.Error()
	}
	if len(out) > 0 {
		r.Lines = splitLines(string(out), lines)
	}
	r.Truncated = len(r.Lines) >= lines
	return r
}

func splitLines(s string, maxLines int) []string {
	raw := strings.Split(strings.TrimRight(s, "\n"), "\n")
	result := make([]string, 0, len(raw))
	for _, l := range raw {
		if l != "" {
			result = append(result, l)
		}
	}
	if len(result) > maxLines {
		result = result[len(result)-maxLines:]
	}
	return result
}
