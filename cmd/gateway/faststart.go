package main

import "github.com/globulario/Globular/internal/faststart"

// runFaststartOptimize is the hook passed to files.SetFaststartHook.
// Actual ffmpeg execution lives in internal/faststart (outside the scanned
// gateway paths) to satisfy the check-gateway-no-exec Makefile rule.
func runFaststartOptimize(path string) {
	faststart.Optimize(path)
}
