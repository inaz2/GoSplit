//go:build !windows

package gosplit

import (
	"golang.org/x/sys/unix"
)

// getDiskFreeSpace returns free disk space where dirPath exists.
func getDiskFreeSpace(dirPath string) (uint64, error) {
	var stat unix.Statfs_t

	if err := unix.Statfs(dirPath, &stat); err != nil {
		return 0, err
	}
	freeBytesAvailable := stat.Bavail * uint64(stat.Bsize)
	return freeBytesAvailable, nil
}
