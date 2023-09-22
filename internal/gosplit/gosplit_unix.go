//go:build !windows

package gosplit

import (
	. "inaz2/GoSplit/internal/gerrors"

	"golang.org/x/sys/unix"
)

// getDiskFreeSpace returns free disk space where dirPath exists.
func getDiskFreeSpace(dirPath string) (uint64, Gerror) {
	var stat unix.Statfs_t

	if err := unix.Statfs(dirPath, &stat); err != nil {
		return 0, GoSplitErrorf("faied to unix.Statfs: %w", err)
	}
	freeBytesAvailable := stat.Bavail * uint64(stat.Bsize)
	return freeBytesAvailable, nil
}
