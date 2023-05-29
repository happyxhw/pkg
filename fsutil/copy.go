package fsutil

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copy file
func CopyFile(src, dst string) (int64, error) {
	srcFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !srcFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	n, err := io.Copy(destFile, srcFile)

	return n, err
}
