//go:build windows

package gosplit

import (
	. "inaz2/GoSplit/internal/gerrors"

	"path/filepath"

	"golang.org/x/sys/windows"
)

// getDiskFreeSpace returns free disk space where dirPath exists.
func getDiskFreeSpace(dirPath string) (uint64, Gerror) {
	var (
		freeBytesAvailableToCaller uint64
		totalNumberOfBytes         uint64
		totalNumberOfFreeBytes     uint64
	)

	dirPath = filepath.FromSlash(dirPath)
	err := windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(dirPath), &freeBytesAvailableToCaller, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, GoSplitErrorf("failed to windows.GetFreeSpaceEx: %w", err)
	}
	return freeBytesAvailableToCaller, nil
}
