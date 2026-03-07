//go:build linux

package stats

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// collect gathers all host + Go runtime metrics and builds a StatsResponse.
func collect(tp TimeProvider) StatsResponse {
	hostname, _ := os.Hostname()
	cpu := collectCPU()
	mem := collectMemory()
	disk := collectDisk("/")
	net := collectNetwork()
	goStats := collectGo()

	return StatsResponse{
		Hostname:  hostname,
		UptimeSec: time.Since(tp.StartTime()).Seconds(),
		CPU:       cpu,
		Memory:    mem,
		Disk:      disk,
		Network:   net,
		Go:        goStats,
	}
}

// ── CPU (/proc/stat) ────────────────────────────────────────────────────────

type cpuSample struct {
	idle  uint64
	total uint64
}

func readCPUSamples() (overall cpuSample, perCore []cpuSample) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}
			var vals [10]uint64
			for i := 1; i < len(fields) && i <= 10; i++ {
				vals[i-1], _ = strconv.ParseUint(fields[i], 10, 64)
			}
			// user, nice, system, idle, iowait, irq, softirq, steal
			idle := vals[3] + vals[4]
			var total uint64
			for i := 0; i < len(fields)-1 && i < 10; i++ {
				total += vals[i]
			}
			s := cpuSample{idle: idle, total: total}
			if fields[0] == "cpu" {
				overall = s
			} else {
				perCore = append(perCore, s)
			}
		}
	}
	return
}

func collectCPU() CPUStats {
	ov1, pc1 := readCPUSamples()
	time.Sleep(200 * time.Millisecond)
	ov2, pc2 := readCPUSamples()

	pct := func(a, b cpuSample) float64 {
		dt := float64(b.total - a.total)
		if dt == 0 {
			return 0
		}
		di := float64(b.idle - a.idle)
		return (1 - di/dt) * 100
	}

	perCore := make([]float64, len(pc1))
	for i := range pc1 {
		if i < len(pc2) {
			perCore[i] = pct(pc1[i], pc2[i])
		}
	}
	return CPUStats{
		Count:    runtime.NumCPU(),
		UsagePct: pct(ov1, ov2),
		PerCore:  perCore,
	}
}

// ── Memory (/proc/meminfo) ──────────────────────────────────────────────────

func collectMemory() MemStats {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemStats{}
	}
	defer f.Close()

	var total, available uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "MemTotal:"):
			total = parseKB(line)
		case strings.HasPrefix(line, "MemAvailable:"):
			available = parseKB(line)
		}
	}
	used := total - available
	var pct float64
	if total > 0 {
		pct = float64(used) / float64(total) * 100
	}
	return MemStats{
		TotalBytes: total,
		UsedBytes:  used,
		UsedPct:    pct,
	}
}

func parseKB(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	v, _ := strconv.ParseUint(fields[1], 10, 64)
	return v * 1024 // kB → bytes
}

// ── Disk (syscall.Statfs) ───────────────────────────────────────────────────

func collectDisk(path string) DiskStats {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskStats{Path: path}
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	var freePct float64
	if total > 0 {
		freePct = float64(free) / float64(total) * 100
	}
	return DiskStats{
		TotalBytes: total,
		UsedBytes:  used,
		FreePct:    freePct,
		Path:       path,
	}
}

// ── Network (/proc/net/dev) ─────────────────────────────────────────────────

func collectNetwork() NetStats {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return NetStats{}
	}
	defer f.Close()

	var rxTotal, txTotal uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		// Skip header lines (contain "|")
		if strings.Contains(line, "|") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		iface := strings.TrimSuffix(fields[0], ":")
		if iface == "lo" {
			continue
		}
		rx, _ := strconv.ParseUint(fields[1], 10, 64)
		tx, _ := strconv.ParseUint(fields[9], 10, 64)
		rxTotal += rx
		txTotal += tx
	}
	return NetStats{RxBytes: rxTotal, TxBytes: txTotal}
}

// ── Go runtime ──────────────────────────────────────────────────────────────

func collectGo() GoStats {
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
