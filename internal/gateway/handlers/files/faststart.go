package files

// faststart.go — background MP4 faststart scheduling hook.
//
// The gateway detects MP4 files whose moov atom is not at the front and calls
// an injected OptimizeFunc hook for each such file.  The hook is set once at
// startup by the binary (cmd/globular-gateway) which wires in the actual
// ffmpeg runner.  The gateway package itself contains no process execution.
//
// Eligible types: .mp4 .m4v .mov .m4a

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
	"sync"
)

// OptimizeFunc is the callback signature for the faststart optimization hook.
// It receives the absolute path of the file to be optimized and runs
// asynchronously (called from its own goroutine).
type OptimizeFunc func(path string)

// faststartVideoExts lists the file extensions eligible for moov-atom
// relocation.  Keep the list tight — only ISO base media file format
// containers benefit from -movflags faststart.
var faststartVideoExts = []string{".mp4", ".m4v", ".mov", ".m4a"}

func isFaststartEligible(name string) bool {
	lname := strings.ToLower(name)
	for _, ext := range faststartVideoExts {
		if strings.HasSuffix(lname, ext) {
			return true
		}
	}
	return false
}

// isMp4FastStart reports whether the MP4 file at path already has its moov
// atom positioned before any mdat box (i.e. it is already "faststart").
//
// It reads only the minimum bytes needed to walk the top-level box headers —
// typically under 100 bytes — so it is safe to call on NFS-mounted files.
//
// Returns true  → moov comes first, skip optimization.
// Returns false → moov is at the end (or file is unreadable), try to optimize.
func isMp4FastStart(path string) bool {
	f, err := os.Open(path) // #nosec G304 — path already sanitized by caller
	if err != nil {
		return false
	}
	defer f.Close()

	// Walk up to 20 top-level boxes looking for moov or mdat.
	// Box layout (ISO 14496-12):
	//   [4 bytes: size][4 bytes: type][...data...]
	// Special sizes:
	//   size == 0  →  box extends to end of file
	//   size == 1  →  actual size in the next 8 bytes (64-bit extended size)
	for i := 0; i < 20; i++ {
		var hdr [8]byte
		if _, err := io.ReadFull(f, hdr[:]); err != nil {
			return false
		}

		size := binary.BigEndian.Uint32(hdr[:4])
		boxType := string(hdr[4:8])

		switch boxType {
		case "moov":
			return true // moov before mdat → already faststart
		case "mdat":
			return false // mdat before moov → not faststart
		}

		// Seek past the body of this box to reach the next one.
		var skip int64
		switch size {
		case 0:
			// Rest-of-file box; moov must have appeared already.
			return true
		case 1:
			// 64-bit extended size: the next 8 bytes hold the true length.
			var ext [8]byte
			if _, err := io.ReadFull(f, ext[:]); err != nil {
				return false
			}
			realSize := int64(binary.BigEndian.Uint64(ext[:]))
			skip = realSize - 16 // already consumed: 8-byte hdr + 8-byte ext
		default:
			skip = int64(size) - 8 // already consumed the 8-byte header
		}

		if skip < 0 {
			return false // malformed box
		}
		if _, err := f.Seek(skip, io.SeekCurrent); err != nil {
			return false
		}
	}

	// 20 boxes without moov or mdat — unusual; assume already optimized to
	// avoid scheduling the same file repeatedly.
	return true
}

// faststartOptimizer deduplicates background optimization jobs.
type faststartOptimizer struct {
	mu         sync.Mutex
	inProgress map[string]struct{}
	hook       OptimizeFunc // injected by the binary; nil → feature disabled
}

// globalFaststartOptimizer is the process-wide singleton.
var globalFaststartOptimizer = &faststartOptimizer{
	inProgress: make(map[string]struct{}),
}

// SetFaststartHook registers the OptimizeFunc that will be called for each
// file that needs moov relocation.  Call this once during server startup from
// the binary (cmd/globular-gateway).  Passing nil disables the feature.
func SetFaststartHook(fn OptimizeFunc) {
	globalFaststartOptimizer.mu.Lock()
	globalFaststartOptimizer.hook = fn
	globalFaststartOptimizer.mu.Unlock()
}

// Schedule enqueues a background optimization for path when all of the
// following conditions are met:
//   - a hook has been registered via SetFaststartHook
//   - the file extension is eligible (mp4, m4v, mov, m4a)
//   - isMp4FastStart returns false (moov is not at the front)
//   - no optimization is already in progress for this path
//
// Returns immediately; work happens in a goroutine.
func (o *faststartOptimizer) Schedule(path string) {
	o.mu.Lock()
	hook := o.hook
	o.mu.Unlock()

	if hook == nil {
		return // feature disabled — no hook registered
	}
	if !isFaststartEligible(path) {
		return
	}
	if isMp4FastStart(path) {
		return // already optimized
	}

	o.mu.Lock()
	if _, running := o.inProgress[path]; running {
		o.mu.Unlock()
		return // already in progress
	}
	o.inProgress[path] = struct{}{}
	o.mu.Unlock()

	go func() {
		defer func() {
			o.mu.Lock()
			delete(o.inProgress, path)
			o.mu.Unlock()
		}()
		hook(path)
	}()
}
