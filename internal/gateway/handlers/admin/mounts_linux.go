//go:build linux

package admin

import (
	"bufio"
	"os"
	"strings"
	"syscall"
)

// skipFSTypes are virtual/pseudo filesystems that should be excluded.
var skipFSTypes = map[string]bool{
	"tmpfs":         true,
	"sysfs":         true,
	"proc":          true,
	"devpts":        true,
	"devtmpfs":      true,
	"cgroup":        true,
	"cgroup2":       true,
	"overlay":       true,
	"squashfs":      true,
	"securityfs":    true,
	"debugfs":       true,
	"tracefs":       true,
	"hugetlbfs":     true,
	"mqueue":        true,
	"pstore":        true,
	"binfmt_misc":   true,
	"configfs":      true,
	"fusectl":       true,
	"autofs":        true,
	"efivarfs":      true,
	"bpf":           true,
	"nsfs":          true,
	"ramfs":         true,
	"rpc_pipefs":    true,
	"nfsd":          true,
	"fuse.snapfuse": true,
	"fuse.lxcfs":    true,
}

// collectMounts reads /proc/mounts, filters to real block devices, and
// calls syscall.Statfs for each to get capacity information.
func collectMounts() []MountInfo {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil
	}
	defer f.Close()

	seen := make(map[string]bool)
	var mounts []MountInfo

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}
		device := fields[0]
		mountPoint := fields[1]
		fsType := fields[2]

		// Skip virtual filesystems
		if skipFSTypes[fsType] {
			continue
		}

		// Skip if not a real block device (must start with /)
		if !strings.HasPrefix(device, "/") {
			continue
		}

		// Skip snap mounts
		if strings.HasPrefix(mountPoint, "/snap/") {
			continue
		}

		// Deduplicate by mount point
		if seen[mountPoint] {
			continue
		}
		seen[mountPoint] = true

		mi := MountInfo{
			Device:     device,
			MountPoint: mountPoint,
			FSType:     fsType,
		}

		// Get filesystem stats
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountPoint, &stat); err == nil {
			mi.TotalBytes = stat.Blocks * uint64(stat.Bsize)
			mi.FreeBytes = stat.Bavail * uint64(stat.Bsize) // available to unprivileged user
			if mi.TotalBytes > mi.FreeBytes {
				mi.UsedBytes = mi.TotalBytes - mi.FreeBytes
			}
		}

		mounts = append(mounts, mi)
	}

	return mounts
}
