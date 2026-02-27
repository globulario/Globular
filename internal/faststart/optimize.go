// Package faststart provides an ffmpeg-based MP4 moov-atom optimizer.
//
// It is intentionally kept in its own package (outside internal/gateway and
// cmd/globular-gateway) so that os/exec is never imported by the gateway
// package itself, satisfying the check-gateway-no-exec Makefile guard.
package faststart

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

// Optimize runs `ffmpeg -movflags faststart -codec copy` on path and
// atomically replaces the original file on success.
//
// It is designed to be called from a goroutine (e.g. via files.SetFaststartHook)
// and never returns an error — failures are logged and the original file is
// left untouched.
func Optimize(path string) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		// ffmpeg not installed — silently disable.
		return
	}

	tmp := path + ".faststart.tmp"

	// Remove any leftover temp file from a previous interrupted run.
	_ = os.Remove(tmp)

	slog.Info("faststart: starting", "file", filepath.Base(path))

	cmd := exec.Command(ffmpeg,
		"-i", path,
		"-movflags", "faststart",
		"-codec", "copy", // remux only — no re-encoding, fast even over NFS
		"-y", // overwrite output without interactive prompt
		tmp,
	)
	// Suppress ffmpeg's verbose output so it doesn't pollute the gateway log.
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		slog.Warn("faststart: ffmpeg failed", "file", filepath.Base(path), "err", err)
		_ = os.Remove(tmp)
		return
	}

	// Sanity-check: output must be at least 4 KB to be a plausible MP4.
	fi, err := os.Stat(tmp)
	if err != nil || fi.Size() < 4096 {
		slog.Warn("faststart: output too small, discarding", "file", filepath.Base(path))
		_ = os.Remove(tmp)
		return
	}

	// Atomic replace.  Any in-flight http.ServeFile holds its own fd and is
	// unaffected; the next request opens the newly optimized file.
	if err := os.Rename(tmp, path); err != nil {
		slog.Warn("faststart: rename failed", "file", filepath.Base(path), "err", err)
		_ = os.Remove(tmp)
		return
	}

	slog.Info("faststart: done", "file", filepath.Base(path))
}
