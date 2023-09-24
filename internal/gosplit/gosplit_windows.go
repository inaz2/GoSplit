//go:build windows

package gosplit

import (
	g "inaz2/GoSplit/internal/gerrors"

	"path/filepath"

	"golang.org/x/sys/windows"
)

// getDiskFreeSpace returns free disk space where dirPath exists.
func getDiskFreeSpace(dirPath string) (uint64, g.Error) {
	var (
		freeBytesAvailableToCaller uint64
		totalNumberOfBytes         uint64
		totalNumberOfFreeBytes     uint64
	)

	dirPath = filepath.FromSlash(dirPath)
	err := windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(dirPath), &freeBytesAvailableToCaller, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, wrapper.Errorf("failed to windows.GetFreeSpaceEx: %w", err)
	}
	return freeBytesAvailableToCaller, nil
}
