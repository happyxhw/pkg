//go:build windows
// +build windows

package fsutil

import (
	"path/filepath"
	"syscall"
)

const dotCharacter = 46

func isHidden(path string) bool {
	// dotfiles also count as hidden (if you want)
	if path[0] == dotCharacter {
		return true
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Appending `\\?\` to the absolute path helps with
	// preventing 'Path Not Specified Error' when accessing
	// long paths and filenames
	// https://docs.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation?tabs=cmd
	pointer, err := syscall.UTF16PtrFromString(`\\?\` + absPath)
	if err != nil {
		return false
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false
	}

	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
