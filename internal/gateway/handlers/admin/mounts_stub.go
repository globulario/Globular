//go:build !linux

package admin

// collectMounts returns a single root mount with zero capacity values
// on non-Linux platforms (macOS, Windows).
func collectMounts() []MountInfo {
	return []MountInfo{{
		Device:     "/dev/stub",
		MountPoint: "/",
		FSType:     "stub",
	}}
}
