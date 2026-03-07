//go:build !linux

package stats

import (
	"os"
	"runtime"
	"time"
)

// collect returns zero-value metrics on non-Linux platforms (macOS dev compat).
func collect(tp TimeProvider) StatsResponse {
	hostname, _ := os.Hostname()
	return StatsResponse{
		Hostname:  hostname,
		UptimeSec: time.Since(tp.StartTime()).Seconds(),
		CPU: CPUStats{
			Count:   runtime.NumCPU(),
			PerCore: make([]float64, runtime.NumCPU()),
		},
		Go: collectGoStub(),
	}
}

func collectGoStub() GoStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var lastPause uint64
	if m.NumGC > 0 {
		idx := (m.NumGC + 255) % 256
		lastPause = m.PauseNs[idx]
	}
	return GoStats{
		Goroutines: runtime.NumGoroutine(),
		HeapAlloc:  m.HeapAlloc,
		GCPauseNs:  lastPause,
		NumGC:      m.NumGC,
	}
}
