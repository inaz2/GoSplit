//go:build windows

package gosplit

import (
	"path/filepath"

	"golang.org/x/sys/windows"
)

// getDiskFreeSpace returns free disk space where dirPath exists.
func getDiskFreeSpace(dirPath string) (uint64, error) {
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
