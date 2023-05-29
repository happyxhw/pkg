//go:build !windows
// +build !windows

package fsutil

const dotCharacter = 46

func isHidden(path string) bool {
	return path[0] == dotCharacter
}
